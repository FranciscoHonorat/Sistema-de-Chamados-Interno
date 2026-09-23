package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/valueobjects"
)

type AssignTicketInput struct {
	Actor      actor.Actor
	TicketID   string
	AssigneeID string
}

type AssignTicketUseCase struct {
	application.EventSourcedUseCase
	responsibles out.ResponsibleDirectory
}

func NewAssignTicketUseCase(store out.EventStore, cache out.TicketCache, responsibles out.ResponsibleDirectory) *AssignTicketUseCase {
	return &AssignTicketUseCase{
		EventSourcedUseCase: application.NewEventSourcedUseCase(store, cache),
		responsibles:        responsibles,
	}
}

func ensureResponsible(ctx context.Context, responsibles out.ResponsibleDirectory, id string) error {
	exists, err := responsibles.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domainErr.ErrUnknownAssignee
	}
	return nil
}

func (uc *AssignTicketUseCase) Execute(ctx context.Context, input AssignTicketInput) error {
	return uc.UpdateTicket(ctx, input.Actor, input.TicketID, ticket.CanManage, func(t *ticket.Ticket) error {
		assigneeID, err := valueobjects.NewAssigneeID(input.AssigneeID)
		if err != nil {
			return err
		}
		if err := ensureResponsible(ctx, uc.responsibles, input.AssigneeID); err != nil {
			return err
		}

		return t.AssignTo(&assigneeID)
	})
}
