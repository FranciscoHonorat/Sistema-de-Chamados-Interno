package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/valueobjects"
)

type OpenTicketInput struct {
	Actor       actor.Actor
	Title       string
	Description string
	AssigneeID  string
	AutoAssign  bool
	Priority    string
}

type OpenTicketOutput struct {
	TicketID string `json:"ticket_id"`
}

type OpenTicketUseCase struct {
	application.EventSourcedUseCase
	responsibles out.ResponsibleDirectory
}

func NewOpenTicketUseCase(store out.EventStore, cache out.TicketCache, responsibles out.ResponsibleDirectory) *OpenTicketUseCase {
	return &OpenTicketUseCase{
		EventSourcedUseCase: application.NewEventSourcedUseCase(store, cache),
		responsibles:        responsibles,
	}
}

func (uc *OpenTicketUseCase) Execute(ctx context.Context, input OpenTicketInput) (OpenTicketOutput, error) {
	title, err := valueobjects.NewTitle(input.Title)
	if err != nil {
		return OpenTicketOutput{}, err
	}

	description, err := valueobjects.NewDescription(input.Description)
	if err != nil {
		return OpenTicketOutput{}, err
	}

	assigneeID, err := uc.chooseAssignee(ctx, input)
	if err != nil {
		return OpenTicketOutput{}, err
	}

	rawPriority := input.Priority
	if rawPriority == "" {
		rawPriority = string(valueobjects.TicketPriorityMedium)
	}
	priority, err := valueobjects.NewPriority(rawPriority)
	if err != nil {
		return OpenTicketOutput{}, err
	}

	id := valueobjects.NewID(uuid.Nil)

	t, err := ticket.NewTicket(id, title, description, valueobjects.TicketStatusOpen, assigneeID, &priority, input.Actor.ID())
	if err != nil {
		return OpenTicketOutput{}, err
	}

	if err := uc.CreateTicket(ctx, input.Actor, t); err != nil {
		return OpenTicketOutput{}, err
	}

	return OpenTicketOutput{TicketID: t.GetID().String()}, nil
}

func (uc *OpenTicketUseCase) chooseAssignee(ctx context.Context, input OpenTicketInput) (*valueobjects.AssigneeID, error) {
	switch {
	case input.AutoAssign && input.AssigneeID != "":
		return nil, domainErr.ErrInvalidAssignee
	case input.AutoAssign:
		return pickLeastBusyResponsible(ctx, uc.EventSourcedUseCase, uc.responsibles)
	case input.AssigneeID == "":
		return nil, nil
	}

	if !ticket.CanManage(nil, input.Actor) {
		return nil, domainErr.ErrForbidden
	}
	assigneeID, err := valueobjects.NewAssigneeID(input.AssigneeID)
	if err != nil {
		return nil, err
	}
	if err := ensureResponsible(ctx, uc.responsibles, input.AssigneeID); err != nil {
		return nil, err
	}
	return &assigneeID, nil
}
