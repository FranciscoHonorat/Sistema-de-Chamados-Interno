package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application/port/out"
)

type fakeOutboxStore struct {
	pending   []out.OutboxEvent
	published []uuid.UUID
	fetchErr  error
}

func (s *fakeOutboxStore) FetchPending(_ context.Context) ([]out.OutboxEvent, error) {
	return s.pending, s.fetchErr
}

func (s *fakeOutboxStore) MarkPublished(_ context.Context, id uuid.UUID) error {
	s.published = append(s.published, id)
	return nil
}

type fakeEventPublisher struct {
	published []out.OutboxEvent
	err       error
	failOn    uuid.UUID
}

func (p *fakeEventPublisher) Publish(_ context.Context, event out.OutboxEvent) error {
	if p.err != nil {
		return p.err
	}
	if event.ID == p.failOn {
		return assert.AnError
	}
	p.published = append(p.published, event)
	return nil
}

func TestPublishPendingEventsUseCase(t *testing.T) {
	t.Run("should publish pending events and mark them as published", func(t *testing.T) {
		id1, id2 := uuid.New(), uuid.New()
		store := &fakeOutboxStore{pending: []out.OutboxEvent{
			{ID: id1, EventType: "EmployeeRegistered", Payload: []byte(`{"id":"agent-1"}`)},
			{ID: id2, EventType: "EmployeeRegistered", Payload: []byte(`{"id":"agent-2"}`)},
		}}
		publisher := &fakeEventPublisher{}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Len(t, publisher.published, 2)
		assert.Equal(t, []uuid.UUID{id1, id2}, store.published)
	})

	t.Run("should do nothing when there are no pending events", func(t *testing.T) {
		store := &fakeOutboxStore{}
		publisher := &fakeEventPublisher{}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Empty(t, publisher.published)
	})

	t.Run("should not mark as published an event that failed to publish", func(t *testing.T) {
		id1 := uuid.New()
		store := &fakeOutboxStore{pending: []out.OutboxEvent{
			{ID: id1, EventType: "EmployeeRegistered", Payload: []byte(`{}`)},
		}}
		publisher := &fakeEventPublisher{err: assert.AnError}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, store.published)
	})

	t.Run("should keep publishing the rest of the batch after one event fails", func(t *testing.T) {
		id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()
		store := &fakeOutboxStore{pending: []out.OutboxEvent{
			{ID: id1, EventType: "EmployeeRegistered"},
			{ID: id2, EventType: "EmployeeSignedUp"},
			{ID: id3, EventType: "EmployeeRegistered"},
		}}
		publisher := &fakeEventPublisher{failOn: id2}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, id2.String())
		assert.Equal(t, []uuid.UUID{id1, id3}, store.published)
	})

	t.Run("should propagate errors from fetching pending events", func(t *testing.T) {
		store := &fakeOutboxStore{fetchErr: assert.AnError}
		publisher := &fakeEventPublisher{}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
