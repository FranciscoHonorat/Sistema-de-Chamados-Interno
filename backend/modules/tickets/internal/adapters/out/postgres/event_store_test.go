package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/database/dbtest"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/postgres"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	return dbtest.Pool(t, "tickets", postgres.Migrations)
}

func TestEventStore(t *testing.T) {
	pool := testPool(t)
	store := postgres.NewEventStore(pool)

	t.Run("should append and load an event stream", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		opened := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")

		err = store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0, "user-1")
		assert.NoError(t, err)

		loaded, err := store.Load(context.Background(), id.GetID())
		assert.NoError(t, err)
		assert.Len(t, loaded, 1)
		assert.Equal(t, "TicketOpened", loaded[0].EventName())
		assert.Equal(t, id.GetID(), loaded[0].AggregateID())

		reopened, ok := loaded[0].(event.TicketOpened)
		assert.True(t, ok)
		assert.Equal(t, "Valid Title", reopened.Title)
		assert.Equal(t, "user-1", reopened.RequesterID)
	})

	t.Run("should record who performed the change alongside each event", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)
		opened := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")

		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0, "user-1"))
		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{event.NewTicketClosed(id, "Resolvido")}, 1, "agent-1"))

		rows, err := pool.Query(context.Background(),
			`SELECT actor_id FROM ticket_events WHERE aggregate_id = $1 ORDER BY version`, id.GetID())
		require.NoError(t, err)
		defer rows.Close()
		var actors []string
		for rows.Next() {
			var actorID string
			require.NoError(t, rows.Scan(&actorID))
			actors = append(actors, actorID)
		}
		assert.Equal(t, []string{"user-1", "agent-1"}, actors)
	})

	t.Run("should append multiple events across separate calls in order", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		opened := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")
		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0, "user-1"))

		assignee, err := valueobjects.NewAssigneeID("agent-1")
		require.NoError(t, err)
		assigned := event.NewTicketAssigned(id, &assignee)
		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{assigned}, 1, "user-1"))

		loaded, err := store.Load(context.Background(), id.GetID())
		assert.NoError(t, err)
		assert.Len(t, loaded, 2)
		assert.Equal(t, "TicketOpened", loaded[0].EventName())
		assert.Equal(t, "TicketAssigned", loaded[1].EventName())
	})

	t.Run("should return an error when expectedVersion is stale", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		opened := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")
		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0, "user-1"))

		err = store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0, "user-1")
		assert.ErrorIs(t, err, domainErr.ErrConcurrencyConflict)
	})

	t.Run("should return an error when loading a stream that does not exist", func(t *testing.T) {
		_, err := store.Load(context.Background(), uuid.New())
		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should list the aggregate IDs of opened tickets", func(t *testing.T) {
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		id1 := valueobjects.NewID(uuid.New())
		opened1 := event.NewTicketOpened(id1, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")
		require.NoError(t, store.Append(context.Background(), id1.GetID(), []event.Event{opened1}, 0, "user-1"))

		id2 := valueobjects.NewID(uuid.New())
		opened2 := event.NewTicketOpened(id2, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")
		require.NoError(t, store.Append(context.Background(), id2.GetID(), []event.Event{opened2}, 0, "user-1"))

		ids, err := store.ListAggregateIDs(context.Background())
		assert.NoError(t, err)
		assert.Contains(t, ids, id1.GetID())
		assert.Contains(t, ids, id2.GetID())
	})

	t.Run("should load the streams of several tickets in one call", func(t *testing.T) {
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)
		assignee, err := valueobjects.NewAssigneeID("agent-1")
		require.NoError(t, err)

		id1, id2 := valueobjects.NewID(uuid.New()), valueobjects.NewID(uuid.New())
		require.NoError(t, store.Append(context.Background(), id1.GetID(), []event.Event{event.NewTicketOpened(id1, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")}, 0, "user-1"))
		require.NoError(t, store.Append(context.Background(), id1.GetID(), []event.Event{event.NewTicketAssigned(id1, &assignee)}, 1, "admin-1"))
		require.NoError(t, store.Append(context.Background(), id2.GetID(), []event.Event{event.NewTicketOpened(id2, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")}, 0, "user-1"))
		unknown := uuid.New()

		histories, err := store.LoadMany(context.Background(), []uuid.UUID{id1.GetID(), id2.GetID(), unknown})

		require.NoError(t, err)
		require.Len(t, histories[id1.GetID()], 2)
		assert.Equal(t, "TicketOpened", histories[id1.GetID()][0].EventName())
		assert.Equal(t, "TicketAssigned", histories[id1.GetID()][1].EventName())
		assert.Len(t, histories[id2.GetID()], 1)
		assert.NotContains(t, histories, unknown)
	})

	t.Run("should load nothing for no tickets", func(t *testing.T) {
		histories, err := store.LoadMany(context.Background(), nil)

		require.NoError(t, err)
		assert.Empty(t, histories)
	})
}
