package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/query"
)

func TestListTicketsUseCase(t *testing.T) {
	t.Run("should return an empty list when there are no tickets", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := query.NewListTicketsUseCase(store, newTestCache())

		output, err := uc.Execute(context.Background(), testUser)

		assert.NoError(t, err)
		assert.Empty(t, output)
	})

	t.Run("should list tickets in the order they were opened", func(t *testing.T) {
		store := outtest.NewEventStore()
		openUC := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())

		first, err := openUC.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "First Ticket",
			Description: "First Description",
		})
		require.NoError(t, err)

		second, err := openUC.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Second Ticket",
			Description: "Second Description",
		})
		require.NoError(t, err)

		uc := query.NewListTicketsUseCase(store, newTestCache())
		output, err := uc.Execute(context.Background(), testUser)

		require.NoError(t, err)
		require.Len(t, output, 2)
		assert.Equal(t, first.TicketID, output[0].TicketID)
		assert.Equal(t, "First Ticket", output[0].Title)
		assert.Equal(t, second.TicketID, output[1].TicketID)
		assert.Equal(t, "Second Ticket", output[1].Title)
	})

	t.Run("should rebuild every uncached ticket from a single read of the store", func(t *testing.T) {
		store := outtest.NewEventStore()
		for _, title := range []string{"First", "Second", "Third"} {
			_, err := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents()).Execute(context.Background(), command.OpenTicketInput{
				Actor: testUser, Title: title, Description: "Description",
			})
			require.NoError(t, err)
		}

		output, err := query.NewListTicketsUseCase(store, cache.NoopTicketCache{}).Execute(context.Background(), testUser)

		require.NoError(t, err)
		require.Len(t, output, 3)
		assert.Equal(t, []string{"First", "Second", "Third"}, []string{output[0].Title, output[1].Title, output[2].Title})
		assert.Equal(t, 1, store.LoadManyCalls)
	})

	t.Run("should not read the store for tickets already in the cache", func(t *testing.T) {
		store := outtest.NewEventStore()
		shared := newTestCache()
		_, err := command.NewOpenTicketUseCase(store, shared, outtest.Agents()).Execute(context.Background(), command.OpenTicketInput{
			Actor: testUser, Title: "Cached", Description: "Description",
		})
		require.NoError(t, err)

		output, err := query.NewListTicketsUseCase(store, shared).Execute(context.Background(), testUser)

		require.NoError(t, err)
		require.Len(t, output, 1)
		assert.Zero(t, store.LoadManyCalls)
	})

	t.Run("should reflect the current state of a mutated ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		opened, err := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents()).Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Valid Title",
			Description: "Valid Description",
		})
		require.NoError(t, err)

		agent := outtest.Actor("agent-1", "support")
		require.NoError(t, command.NewAssignTicketUseCase(store, newTestCache(), outtest.Agents()).Execute(context.Background(), command.AssignTicketInput{
			Actor: testAdmin, TicketID: opened.TicketID, AssigneeID: agent.ID(),
		}))
		err = command.NewCloseTicketUseCase(store, newTestCache()).Execute(context.Background(), command.CloseTicketInput{
			Resolution: "Resolvido",
			Actor:      agent,
			TicketID:   opened.TicketID,
		})
		require.NoError(t, err)

		uc := query.NewListTicketsUseCase(store, newTestCache())
		output, err := uc.Execute(context.Background(), testUser)

		require.NoError(t, err)
		require.Len(t, output, 1)
		assert.Equal(t, "Closed", output[0].Status)
		require.NotNil(t, output[0].ClosedAt)
		assert.Equal(t, "Resolvido", output[0].Resolution)
		assert.False(t, output[0].ClosedAt.Before(output[0].CreatedAt))
	})

	t.Run("should only list the tickets the viewer is allowed to see", func(t *testing.T) {
		store := outtest.NewEventStore()
		openUC := command.NewOpenTicketUseCase(store, newTestCache(), outtest.Agents())
		mine, err := openUC.Execute(context.Background(), command.OpenTicketInput{
			Actor: testUser, Title: "Mine", Description: "Mine",
		})
		require.NoError(t, err)
		_, err = openUC.Execute(context.Background(), command.OpenTicketInput{
			Actor: outtest.Actor("user-2", "user"), Title: "Theirs", Description: "Theirs",
		})
		require.NoError(t, err)
		uc := query.NewListTicketsUseCase(store, newTestCache())

		forUser, err := uc.Execute(context.Background(), testUser)
		require.NoError(t, err)
		forAdmin, err := uc.Execute(context.Background(), outtest.Actor("admin-1", "admin"))
		require.NoError(t, err)

		require.Len(t, forUser, 1)
		assert.Equal(t, mine.TicketID, forUser[0].TicketID)
		assert.Len(t, forAdmin, 2)
	})
}
