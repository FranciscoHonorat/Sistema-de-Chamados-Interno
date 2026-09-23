package out

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
)

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (actor.Actor, error)
}
