package eventbus

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application/port/out"
)

var _ out.EventPublisher = (*Publisher)(nil)

// Publisher hands the events the outbox relay picks up to the in-process bus,
// where the other modules subscribe to them.
type Publisher struct {
	bus *eventbus.Bus
}

func NewPublisher(bus *eventbus.Bus) *Publisher {
	return &Publisher{bus: bus}
}

func (p *Publisher) Publish(ctx context.Context, event out.OutboxEvent) error {
	return p.bus.Publish(ctx, eventbus.Message{ID: event.ID.String(), Type: event.EventType, Payload: event.Payload})
}
