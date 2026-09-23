package postgres

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application/port/out"
)

const (
	outboxBatchSize = 100
	outboxLease     = 15 * time.Second
)

type OutboxStore struct {
	pool *pgxpool.Pool
}

func NewOutboxStore(pool *pgxpool.Pool) *OutboxStore {
	return &OutboxStore{pool: pool}
}

// FetchPending claims the oldest unpublished events that no other relay holds,
// leasing them for a few seconds. They come back in the order they happened.
func (s *OutboxStore) FetchPending(ctx context.Context) ([]out.OutboxEvent, error) {
	rows, err := s.pool.Query(ctx,
		`UPDATE outbox_events SET locked_until = now() + make_interval(secs => $2)
		 WHERE id IN (
		     SELECT id FROM outbox_events
		     WHERE published_at IS NULL AND (locked_until IS NULL OR locked_until < now())
		     ORDER BY occurred_at ASC
		     LIMIT $1
		     FOR UPDATE SKIP LOCKED
		 )
		 RETURNING id, event_type, payload, occurred_at`,
		outboxBatchSize, outboxLease.Seconds(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type claimed struct {
		event      out.OutboxEvent
		occurredAt time.Time
	}
	var batch []claimed
	for rows.Next() {
		var c claimed
		if err := rows.Scan(&c.event.ID, &c.event.EventType, &c.event.Payload, &c.occurredAt); err != nil {
			return nil, err
		}
		batch = append(batch, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.SliceStable(batch, func(i, j int) bool { return batch[i].occurredAt.Before(batch[j].occurredAt) })
	events := make([]out.OutboxEvent, len(batch))
	for i, c := range batch {
		events[i] = c.event
	}
	return events, nil
}

func (s *OutboxStore) MarkPublished(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `UPDATE outbox_events SET published_at = now(), locked_until = NULL WHERE id = $1`, id)
	return err
}
