package out

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/event"
)

type EventStore interface {
	Append(ctx context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int, actorID string) error
	Load(ctx context.Context, aggregateID uuid.UUID) ([]event.Event, error)
	LoadMany(ctx context.Context, aggregateIDs []uuid.UUID) (map[uuid.UUID][]event.Event, error)
	ListAggregateIDs(ctx context.Context) ([]uuid.UUID, error)
}
