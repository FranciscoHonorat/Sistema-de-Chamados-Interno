// Package events subscribes the tickets module to the integration events the
// employees module publishes on the in-process bus.
package events

import (
	"context"
	"encoding/json"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/contracts"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/command"
)

type responsibleSyncer interface {
	Execute(ctx context.Context, input command.SyncResponsibleInput) error
}

type accountNotifier interface {
	Execute(ctx context.Context, input command.AccountRequestInput) error
}

type EmployeesSubscriber struct {
	syncResponsible responsibleSyncer
	notifyAccount   accountNotifier
}

func NewEmployeesSubscriber(syncResponsible responsibleSyncer, notifyAccount accountNotifier) *EmployeesSubscriber {
	return &EmployeesSubscriber{syncResponsible: syncResponsible, notifyAccount: notifyAccount}
}

// Subscribe registers the subscriber for every event type it handles.
func (s *EmployeesSubscriber) Subscribe(bus *eventbus.Bus) {
	bus.Subscribe(contracts.EventEmployeeRegistered, "tickets.sync-responsible", s.Handle)
	bus.Subscribe(contracts.EventEmployeeSignedUp, "tickets.notify-account-request", s.Handle)
	bus.Subscribe(contracts.EventPasswordResetRequested, "tickets.notify-account-request", s.Handle)
}

var accountRequests = map[string]command.AccountRequest{
	contracts.EventEmployeeSignedUp:       command.AccountSignUp,
	contracts.EventPasswordResetRequested: command.AccountPasswordReset,
}

func (s *EmployeesSubscriber) Handle(ctx context.Context, msg eventbus.Message) error {
	if msg.Type == contracts.EventEmployeeRegistered {
		var payload contracts.EmployeeRegistered
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}
		return s.syncResponsible.Execute(ctx, command.SyncResponsibleInput{ID: payload.ID, Name: payload.Name, Role: payload.Role})
	}

	request, isAccountRequest := accountRequests[msg.Type]
	if !isAccountRequest {
		return nil
	}
	// EmployeeSignedUp and PasswordResetRequested share the same shape.
	var payload contracts.EmployeeSignedUp
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	return s.notifyAccount.Execute(ctx, command.AccountRequestInput{EventID: msg.ID, Kind: request, EmployeeID: payload.ID, Name: payload.Name})
}
