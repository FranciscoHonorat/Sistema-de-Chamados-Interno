package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/notification"
)

var accountRequestNamespace = uuid.MustParse("5b0f2a4e-8f3c-4d7a-9c61-2e8b7d4a1f90")

type AccountRequest string

const (
	AccountSignUp        AccountRequest = "signup"
	AccountPasswordReset AccountRequest = "password_reset"
)

var accountNotifications = map[AccountRequest]func(employeeID, name string, at time.Time) notification.Notification{
	AccountSignUp:        notification.ForNewAccount,
	AccountPasswordReset: notification.ForPasswordResetRequest,
}

type AccountRequestInput struct {
	EventID    string
	Kind       AccountRequest
	EmployeeID string
	Name       string
}

type NotifyAccountRequestUseCase struct {
	notifications out.NotificationStore
	now           func() time.Time
}

func NewNotifyAccountRequestUseCase(notifications out.NotificationStore, now func() time.Time) *NotifyAccountRequestUseCase {
	return &NotifyAccountRequestUseCase{notifications: notifications, now: now}
}

func (uc *NotifyAccountRequestUseCase) Execute(ctx context.Context, input AccountRequestInput) error {
	notify, known := accountNotifications[input.Kind]
	if !known {
		return nil
	}
	n := notify(input.EmployeeID, input.Name, uc.now())
	if input.EventID != "" {
		n.ID = uuid.NewSHA1(accountRequestNamespace, []byte(input.EventID))
	}
	return uc.notifications.Add(ctx, []notification.Notification{n})
}
