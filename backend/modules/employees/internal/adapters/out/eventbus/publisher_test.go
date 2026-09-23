package eventbus_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bus "github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/adapters/out/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application/port/out"
)

func TestPublisher(t *testing.T) {
	t.Run("should hand the event to the subscribers of its type", func(t *testing.T) {
		b := bus.New()
		var received bus.Message
		b.Subscribe("EmployeeRegistered", "test", func(_ context.Context, m bus.Message) error {
			received = m
			return nil
		})

		id := uuid.New()
		err := eventbus.NewPublisher(b).Publish(context.Background(), out.OutboxEvent{ID: id, EventType: "EmployeeRegistered", Payload: []byte(`{"id":"agent-1"}`)})

		require.NoError(t, err)
		assert.Equal(t, bus.Message{ID: id.String(), Type: "EmployeeRegistered", Payload: []byte(`{"id":"agent-1"}`)}, received)
	})

	t.Run("should fail when a subscriber fails, so the outbox retries", func(t *testing.T) {
		b := bus.New()
		b.Subscribe("EmployeeRegistered", "test", func(context.Context, bus.Message) error { return assert.AnError })

		err := eventbus.NewPublisher(b).Publish(context.Background(), out.OutboxEvent{ID: uuid.New(), EventType: "EmployeeRegistered"})

		assert.ErrorIs(t, err, assert.AnError)
	})
}
