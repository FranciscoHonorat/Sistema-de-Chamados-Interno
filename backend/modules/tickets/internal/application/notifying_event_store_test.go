package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/notification"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/valueobjects"
)

func openedEvent(t *testing.T) (uuid.UUID, event.Event) {
	t.Helper()
	id := valueobjects.NewID(uuid.New())
	title, err := valueobjects.NewTitle("Impressora")
	require.NoError(t, err)
	description, err := valueobjects.NewDescription("Não imprime")
	require.NoError(t, err)
	return id.GetID(), event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")
}

func TestNotifyingEventStore(t *testing.T) {
	t.Run("records the notifications of the events it stores", func(t *testing.T) {
		notifications := &outtest.NotificationStore{}
		store := application.NewNotifyingEventStore(outtest.NewEventStore(), notifications)
		id, opened := openedEvent(t)

		require.NoError(t, store.Append(context.Background(), id, []event.Event{opened}, 0, "user-1"))

		require.Len(t, notifications.Added, 2)
		assert.Equal(t, notification.ToRole("support"), notifications.Added[0].Audience)
		assert.Equal(t, "Novo chamado: Impressora", notifications.Added[0].Message)
	})

	t.Run("notifies nothing when the events could not be stored", func(t *testing.T) {
		notifications := &outtest.NotificationStore{}
		store := application.NewNotifyingEventStore(outtest.NewEventStore(), notifications)
		id, opened := openedEvent(t)

		err := store.Append(context.Background(), id, []event.Event{opened}, 3, "user-1")

		assert.Error(t, err)
		assert.Empty(t, notifications.Added)
	})

	t.Run("keeps the change even if the notifications fail", func(t *testing.T) {
		inner := outtest.NewEventStore()
		store := application.NewNotifyingEventStore(inner, &outtest.NotificationStore{AddErr: assert.AnError})
		id, opened := openedEvent(t)

		require.NoError(t, store.Append(context.Background(), id, []event.Event{opened}, 0, "user-1"))

		stored, err := inner.Load(context.Background(), id)
		require.NoError(t, err)
		assert.Len(t, stored, 1)
	})

	t.Run("builds the notifications from the ticket it was given, without reading the history again", func(t *testing.T) {
		inner := &loadCountingStore{EventStore: outtest.NewEventStore()}
		notifications := &outtest.NotificationStore{}
		store := application.NewNotifyingEventStore(inner, notifications)
		admin := outtest.Actor("admin-1", "admin")

		opened, err := command.NewOpenTicketUseCase(store, cache.NewInMemoryTicketCache(), outtest.Agents()).Execute(context.Background(), command.OpenTicketInput{
			Actor: outtest.Actor("user-1", "user"), Title: "Impressora", Description: "Não imprime",
		})
		require.NoError(t, err)
		require.NoError(t, command.NewAssignTicketUseCase(store, cache.NoopTicketCache{}, outtest.Agents()).Execute(context.Background(), command.AssignTicketInput{
			Actor: admin, TicketID: opened.TicketID, AssigneeID: "agent-1",
		}))

		assert.Equal(t, 1, inner.loads)
		assert.Contains(t, messagesOf(notifications.Added), "Chamado atribuído a você: Impressora")
	})
}

type loadCountingStore struct {
	*outtest.EventStore
	loads int
}

func (s *loadCountingStore) Load(ctx context.Context, aggregateID uuid.UUID) ([]event.Event, error) {
	s.loads++
	return s.EventStore.Load(ctx, aggregateID)
}

func messagesOf(notifications []notification.Notification) []string {
	messages := make([]string, 0, len(notifications))
	for _, n := range notifications {
		messages = append(messages, n.Message)
	}
	return messages
}
