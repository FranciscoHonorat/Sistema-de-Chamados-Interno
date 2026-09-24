package command_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/ticket"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenTicketUseCase(t *testing.T) {
	t.Run("should record the actor as the one who opened the ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		output, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       outtest.Actor("user-7", "user"),
			Title:       "Valid Title",
			Description: "Valid Description",
		})
		assert.NoError(t, err)

		history, err := store.Load(context.Background(), uuid.MustParse(output.TicketID))
		assert.NoError(t, err)
		tk, err := ticket.LoadFromHistory(history)
		assert.NoError(t, err)
		assert.Equal(t, "user-7", tk.GetRequesterID())
	})

	t.Run("should open a new ticket and persist a TicketOpened event", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		output, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testAdmin,
			Title:       "Valid Title",
			Description: "Valid Description",
			AssigneeID:  "agent-1",
			Priority:    "High",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, output.TicketID)

		ticketID, err := uuid.Parse(output.TicketID)
		assert.NoError(t, err)

		history, err := store.Load(context.Background(), ticketID)
		assert.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, "TicketOpened", history[0].EventName())

		tk, err := ticket.LoadFromHistory(history)
		assert.NoError(t, err)
		assert.Equal(t, "Valid Title", tk.GetTitle().GetTitle())
		assert.Equal(t, "agent-1", tk.GetAssigneeID().GetAssigneeID())
	})

	t.Run("should forbid a regular user from choosing the assignee when opening a ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Valid Title",
			Description: "Valid Description",
			AssigneeID:  "agent-1",
		})

		assert.ErrorIs(t, err, domainErr.ErrForbidden)
	})

	t.Run("should refuse to open a ticket assigned to someone who is not a support agent", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testAdmin,
			Title:       "Valid Title",
			Description: "Valid Description",
			AssigneeID:  "user-1",
		})

		assert.ErrorIs(t, err, domainErr.ErrUnknownAssignee)
		ids, err := store.ListAggregateIDs(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("should open a ticket without an assignee and with medium priority by default", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		output, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Valid Title",
			Description: "Valid Description",
		})
		require.NoError(t, err)

		tk := loadOpenedTicket(t, store, output.TicketID)
		assert.Nil(t, tk.GetAssigneeID())
		assert.Equal(t, "Medium", tk.GetPriority().GetPriority())
	})

	t.Run("should let a regular user open a ticket assigned to the least busy responsible", func(t *testing.T) {
		store := outtest.NewEventStore()
		cache := newTestCache()
		responsibles := &outtest.ResponsibleDirectory{Responsibles: []string{"agent-1", "agent-2"}}
		uc := command.NewOpenTicketUseCase(store, cache, responsibles)

		first, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor: testUser, Title: "First", Description: "Valid Description", AutoAssign: true,
		})
		require.NoError(t, err)
		second, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor: testUser, Title: "Second", Description: "Valid Description", AutoAssign: true,
		})
		require.NoError(t, err)

		assert.Equal(t, "agent-1", loadOpenedTicket(t, store, first.TicketID).GetAssigneeID().GetAssigneeID())
		assert.Equal(t, "agent-2", loadOpenedTicket(t, store, second.TicketID).GetAssigneeID().GetAssigneeID())
	})

	t.Run("should refuse automatic assignment when there are no responsibles", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), &outtest.ResponsibleDirectory{})

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor: testUser, Title: "Valid Title", Description: "Valid Description", AutoAssign: true,
		})

		assert.ErrorIs(t, err, domainErr.ErrNoResponsiblesAvailable)
	})

	t.Run("should refuse a manual assignee together with automatic assignment", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor: testAdmin, Title: "Valid Title", Description: "Valid Description", AssigneeID: "agent-1", AutoAssign: true,
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
	})

	t.Run("should return an error for an invalid title", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "",
			Description: "Valid Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error for an invalid priority", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Valid Title",
			Description: "Valid Description",
			Priority:    "bogus",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidPriority)
	})
}

func loadOpenedTicket(t *testing.T, store *outtest.EventStore, ticketID string) *ticket.Ticket {
	t.Helper()
	history, err := store.Load(context.Background(), uuid.MustParse(ticketID))
	require.NoError(t, err)
	tk, err := ticket.LoadFromHistory(history)
	require.NoError(t, err)
	return tk
}
