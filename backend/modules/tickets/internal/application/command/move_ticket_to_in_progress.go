package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/ticket"
)

type MoveTicketToInProgressInput struct {
	Actor    actor.Actor
	TicketID string
}

type MoveTicketToInProgressUseCase struct {
	application.EventSourcedUseCase
}

func NewMoveTicketToInProgressUseCase(store out.EventStore, cache out.TicketCache) *MoveTicketToInProgressUseCase {
	return &MoveTicketToInProgressUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *MoveTicketToInProgressUseCase) Execute(ctx context.Context, input MoveTicketToInProgressInput) error {
	return uc.UpdateTicket(ctx, input.Actor, input.TicketID, ticket.CanWork, func(t *ticket.Ticket) error {
		return t.MoveToInProgress()
	})
}
