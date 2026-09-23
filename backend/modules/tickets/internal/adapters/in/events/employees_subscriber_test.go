package events_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/in/events"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/command"
)

type fakeResponsibleSyncer struct {
	upserted map[string]string
	inputs   []command.SyncResponsibleInput
	err      error
}

func newFakeResponsibleUpserter() *fakeResponsibleSyncer {
	return &fakeResponsibleSyncer{upserted: make(map[string]string)}
}

func (f *fakeResponsibleSyncer) Execute(_ context.Context, input command.SyncResponsibleInput) error {
	if f.err != nil {
		return f.err
	}
	f.upserted[input.ID] = input.Name
	f.inputs = append(f.inputs, input)
	return nil
}

type fakeAccountNotifier struct {
	inputs []command.AccountRequestInput
}

func (f *fakeAccountNotifier) Execute(_ context.Context, input command.AccountRequestInput) error {
	f.inputs = append(f.inputs, input)
	return nil
}

func employeeEvent(eventType, payload string) eventbus.Message {
	return eventbus.Message{Type: eventType, Payload: []byte(payload)}
}

func employeeRegistered(payload string) eventbus.Message {
	return employeeEvent("EmployeeRegistered", payload)
}

func TestEmployeesSubscriber(t *testing.T) {
	t.Run("should upsert a responsible on EmployeeRegistered", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		subscriber := events.NewEmployeesSubscriber(upserter, &fakeAccountNotifier{})

		err := subscriber.Handle(context.Background(), employeeRegistered(`{"id":"agent-1","name":"Ana Souza"}`))

		require.NoError(t, err)
		assert.Equal(t, "Ana Souza", upserter.upserted["agent-1"])
	})

	t.Run("should forward the employee role", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		subscriber := events.NewEmployeesSubscriber(upserter, &fakeAccountNotifier{})

		require.NoError(t, subscriber.Handle(context.Background(), employeeRegistered(`{"id":"admin-1","name":"Admin","role":"admin"}`)))
		assert.Equal(t, []command.SyncResponsibleInput{{ID: "admin-1", Name: "Admin", Role: "admin"}}, upserter.inputs)
	})

	t.Run("should ignore unknown event types", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		subscriber := events.NewEmployeesSubscriber(upserter, &fakeAccountNotifier{})

		err := subscriber.Handle(context.Background(), employeeEvent("SomethingElse", `{}`))

		require.NoError(t, err)
		assert.Empty(t, upserter.upserted)
	})

	t.Run("should return an error for a malformed payload", func(t *testing.T) {
		subscriber := events.NewEmployeesSubscriber(newFakeResponsibleUpserter(), &fakeAccountNotifier{})

		assert.Error(t, subscriber.Handle(context.Background(), employeeRegistered(`not-json`)))
	})

	t.Run("should propagate sync errors", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		upserter.err = assert.AnError
		subscriber := events.NewEmployeesSubscriber(upserter, &fakeAccountNotifier{})

		err := subscriber.Handle(context.Background(), employeeRegistered(`{"id":"agent-1","name":"Ana Souza"}`))

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should tell the admins about a new account", func(t *testing.T) {
		notifier := &fakeAccountNotifier{}
		subscriber := events.NewEmployeesSubscriber(newFakeResponsibleUpserter(), notifier)

		require.NoError(t, subscriber.Handle(context.Background(), employeeEvent("EmployeeSignedUp", `{"id":"user-9","name":"Maria Lima"}`)))

		assert.Equal(t, []command.AccountRequestInput{{Kind: command.AccountSignUp, EmployeeID: "user-9", Name: "Maria Lima"}}, notifier.inputs)
	})

	t.Run("should tell the admins about a password request", func(t *testing.T) {
		notifier := &fakeAccountNotifier{}
		subscriber := events.NewEmployeesSubscriber(newFakeResponsibleUpserter(), notifier)

		require.NoError(t, subscriber.Handle(context.Background(), employeeEvent("PasswordResetRequested", `{"id":"user-1","name":"Usuário Padrão"}`)))

		assert.Equal(t, []command.AccountRequestInput{{Kind: command.AccountPasswordReset, EmployeeID: "user-1", Name: "Usuário Padrão"}}, notifier.inputs)
	})

	t.Run("should pass the event ID on, so a redelivery is recognized", func(t *testing.T) {
		notifier := &fakeAccountNotifier{}
		subscriber := events.NewEmployeesSubscriber(newFakeResponsibleUpserter(), notifier)
		msg := employeeEvent("EmployeeSignedUp", `{"id":"user-9","name":"Maria Lima"}`)
		msg.ID = "event-1"

		require.NoError(t, subscriber.Handle(context.Background(), msg))

		require.Len(t, notifier.inputs, 1)
		assert.Equal(t, "event-1", notifier.inputs[0].EventID)
	})

	t.Run("should receive the events it subscribed to from the bus", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		notifier := &fakeAccountNotifier{}
		bus := eventbus.New()
		events.NewEmployeesSubscriber(upserter, notifier).Subscribe(bus)

		require.NoError(t, bus.Publish(context.Background(), employeeRegistered(`{"id":"agent-1","name":"Ana Souza","role":"support"}`)))
		require.NoError(t, bus.Publish(context.Background(), employeeEvent("EmployeeSignedUp", `{"id":"user-9","name":"Maria Lima"}`)))

		assert.Equal(t, "Ana Souza", upserter.upserted["agent-1"])
		assert.Len(t, notifier.inputs, 1)
	})
}
