package out

import (
	"context"
	"time"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/notification"
)

type NotificationStore interface {
	Add(ctx context.Context, notifications []notification.Notification) error
	ListFor(ctx context.Context, viewer actor.Actor, limit int) ([]notification.Notification, error)
	SeenUntil(ctx context.Context, userID string) (time.Time, error)
	MarkSeen(ctx context.Context, userID string, at time.Time) error
}
