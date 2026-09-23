package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application/port/out"
)

type PublishPendingEventsUseCase struct {
	store     out.OutboxStore
	publisher out.EventPublisher
}

func NewPublishPendingEventsUseCase(store out.OutboxStore, publisher out.EventPublisher) *PublishPendingEventsUseCase {
	return &PublishPendingEventsUseCase{store: store, publisher: publisher}
}

func (uc *PublishPendingEventsUseCase) Execute(ctx context.Context) error {
	events, err := uc.store.FetchPending(ctx)
	if err != nil {
		return err
	}

	var errs []error
	for _, e := range events {
		if err := uc.publisher.Publish(ctx, e); err != nil {
			errs = append(errs, fmt.Errorf("publish %s %s: %w", e.EventType, e.ID, err))
			continue
		}
		if err := uc.store.MarkPublished(ctx, e.ID); err != nil {
			errs = append(errs, fmt.Errorf("mark %s published: %w", e.ID, err))
		}
	}

	return errors.Join(errs...)
}
