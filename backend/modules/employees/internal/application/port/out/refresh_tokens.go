package out

import (
	"context"
	"time"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/domain/session"
)

type RefreshTokenGenerator interface {
	Generate() (string, error)
	Hash(token string) string
}

type RefreshTokenStore interface {
	Save(ctx context.Context, token session.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (session.RefreshToken, bool, error)
	Rotate(ctx context.Context, hash string, at time.Time) (bool, error)
	Revoke(ctx context.Context, hash string) (bool, error)
	RevokeAllForEmployee(ctx context.Context, employeeID string) error
}
