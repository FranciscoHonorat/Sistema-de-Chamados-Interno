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

type AutoAssignTicketInput struct {
	Actor    actor.Actor
	TicketID string
}

type AutoAssignTicketOutput struct {
	AssigneeID string `json:"assignee_id"`
}

type AutoAssignTicketUseCase struct {
	application.EventSourcedUseCase
	responsibles out.ResponsibleDirectory
}

func NewAutoAssignTicketUseCase(store out.EventStore, cache out.TicketCache, responsibles out.ResponsibleDirectory) *AutoAssignTicketUseCase {
	return &AutoAssignTicketUseCase{
		EventSourcedUseCase: application.NewEventSourcedUseCase(store, cache),
		responsibles:        responsibles,
	}
}

func (uc *AutoAssignTicketUseCase) Execute(ctx context.Context, input AutoAssignTicketInput) (AutoAssignTicketOutput, error) {
	var chosenID string
	err := uc.UpdateTicket(ctx, input.Actor, input.TicketID, ticket.CanManage, func(t *ticket.Ticket) error {
		assigneeID, err := pickLeastBusyResponsible(ctx, uc.EventSourcedUseCase, uc.responsibles)
		if err != nil {
			return err
		}
		chosenID = assigneeID.GetAssigneeID()
		return t.AssignTo(assigneeID)
	})
	if err != nil {
		return AutoAssignTicketOutput{}, err
	}

	return AutoAssignTicketOutput{AssigneeID: chosenID}, nil
}

func pickLeastBusyResponsible(ctx context.Context, tickets application.EventSourcedUseCase, responsibles out.ResponsibleDirectory) (*valueobjects.AssigneeID, error) {
	candidates, err := responsibles.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, domainErr.ErrNoResponsiblesAvailable
	}

	allTickets, err := tickets.LoadAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	assigneeID, err := valueobjects.NewAssigneeID(leastBusyResponsible(candidates, countOpenTicketsByAssignee(allTickets)))
	if err != nil {
		return nil, err
	}
	return &assigneeID, nil
}

func countOpenTicketsByAssignee(tickets []*ticket.Ticket) map[string]int {
	counts := make(map[string]int)
	for _, t := range tickets {
		if t.GetStatus() == valueobjects.TicketStatusClosed {
			continue
		}
		if t.GetAssigneeID() == nil {
			continue
		}
		counts[t.GetAssigneeID().GetAssigneeID()]++
	}
	return counts
}

func leastBusyResponsible(candidates []string, openCounts map[string]int) string {
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if openCounts[candidate] < openCounts[best] {
			best = candidate
		}
	}
	return best
}
