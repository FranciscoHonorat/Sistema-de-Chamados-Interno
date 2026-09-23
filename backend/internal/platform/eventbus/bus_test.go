package eventbus_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
)

func TestBus(t *testing.T) {
	t.Run("should deliver a message to every subscriber of its type", func(t *testing.T) {
		bus := eventbus.New()
		var got []string
		bus.Subscribe("A", "first", func(_ context.Context, m eventbus.Message) error {
			got = append(got, "first:"+string(m.Payload))
			return nil
		})
		bus.Subscribe("A", "second", func(_ context.Context, m eventbus.Message) error {
			got = append(got, "second:"+string(m.Payload))
			return nil
		})
		bus.Subscribe("B", "other", func(context.Context, eventbus.Message) error {
			got = append(got, "other")
			return nil
		})

		require.NoError(t, bus.Publish(context.Background(), eventbus.Message{Type: "A", Payload: []byte("x")}))
		assert.Equal(t, []string{"first:x", "second:x"}, got)
	})

	t.Run("should succeed when nobody subscribed", func(t *testing.T) {
		assert.NoError(t, eventbus.New().Publish(context.Background(), eventbus.Message{Type: "A"}))
	})

	t.Run("should keep delivering after a failure and report it", func(t *testing.T) {
		bus := eventbus.New()
		delivered := false
		bus.Subscribe("A", "broken", func(context.Context, eventbus.Message) error { return assert.AnError })
		bus.Subscribe("A", "healthy", func(context.Context, eventbus.Message) error {
			delivered = true
			return nil
		})

		err := bus.Publish(context.Background(), eventbus.Message{Type: "A"})

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "broken")
		assert.True(t, delivered)
	})

	t.Run("should turn a panicking subscriber into an error", func(t *testing.T) {
		bus := eventbus.New()
		bus.Subscribe("A", "panics", func(context.Context, eventbus.Message) error { panic("boom") })

		assert.ErrorContains(t, bus.Publish(context.Background(), eventbus.Message{Type: "A"}), "boom")
	})
}
