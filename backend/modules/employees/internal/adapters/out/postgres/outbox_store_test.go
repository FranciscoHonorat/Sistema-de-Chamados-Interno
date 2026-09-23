package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/domain/employee"
)

func TestOutboxStore(t *testing.T) {
	pool := testPool(t)
	repo := postgres.NewEmployeeRepository(pool)
	store := postgres.NewOutboxStore(pool)
	ctx := context.Background()

	drain := func() {
		for {
			events, err := store.FetchPending(ctx)
			require.NoError(t, err)
			if len(events) == 0 {
				return
			}
			for _, e := range events {
				require.NoError(t, store.MarkPublished(ctx, e.ID))
			}
		}
	}
	drain()

	register := func(id, name string) {
		e, err := employee.Register(id, name, id, employee.RoleSupport, "hash")
		require.NoError(t, err)
		require.NoError(t, repo.Register(ctx, e))
	}

	t.Run("should hand out pending events oldest first", func(t *testing.T) {
		first, second := uniqueID(), uniqueID()
		register(first, "First")
		register(second, "Second")

		events, err := store.FetchPending(ctx)

		require.NoError(t, err)
		require.Len(t, events, 2)
		assert.Contains(t, string(events[0].Payload), first)
		assert.Contains(t, string(events[1].Payload), second)
		for _, e := range events {
			require.NoError(t, store.MarkPublished(ctx, e.ID))
		}
	})

	t.Run("should not hand the same event to two relays while it is leased", func(t *testing.T) {
		id := uniqueID()
		register(id, "Leased")

		claimed, err := store.FetchPending(ctx)
		require.NoError(t, err)
		require.Len(t, claimed, 1)

		again, err := store.FetchPending(ctx)
		require.NoError(t, err)
		assert.Empty(t, again)

		require.NoError(t, store.MarkPublished(ctx, claimed[0].ID))
	})

	t.Run("should not hand out published events again", func(t *testing.T) {
		id := uniqueID()
		register(id, "Published")
		claimed, err := store.FetchPending(ctx)
		require.NoError(t, err)
		require.Len(t, claimed, 1)
		require.NoError(t, store.MarkPublished(ctx, claimed[0].ID))

		_, err = pool.Exec(ctx, `UPDATE outbox_events SET locked_until = NULL WHERE id = $1`, claimed[0].ID)
		require.NoError(t, err)

		events, err := store.FetchPending(ctx)
		require.NoError(t, err)
		for _, e := range events {
			assert.NotEqual(t, claimed[0].ID, e.ID)
		}
	})

	t.Run("should hand out an event again once its lease expires", func(t *testing.T) {
		id := uniqueID()
		register(id, "Expired")
		claimed, err := store.FetchPending(ctx)
		require.NoError(t, err)
		require.Len(t, claimed, 1)

		_, err = pool.Exec(ctx, `UPDATE outbox_events SET locked_until = now() - interval '1 second' WHERE id = $1`, claimed[0].ID)
		require.NoError(t, err)

		events, err := store.FetchPending(ctx)
		require.NoError(t, err)
		require.Len(t, events, 1)
		assert.Equal(t, claimed[0].ID, events[0].ID)
		require.NoError(t, store.MarkPublished(ctx, events[0].ID))
	})

}
