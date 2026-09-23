// Package employees is the tickets module's anti-corruption layer towards the
// employees module: it turns the principal the employees module authenticates
// into the tickets domain's actor.
package employees

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/contracts"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
)

var _ out.TokenVerifier = (*TokenVerifier)(nil)

type TokenVerifier struct {
	identity contracts.TokenVerifier
}

func NewTokenVerifier(identity contracts.TokenVerifier) *TokenVerifier {
	return &TokenVerifier{identity: identity}
}

func (v *TokenVerifier) Verify(ctx context.Context, token string) (actor.Actor, error) {
	principal, err := v.identity.VerifyAccessToken(ctx, token)
	if err != nil {
		return actor.Actor{}, err
	}
	if principal.MustChangePassword {
		return actor.Actor{}, domainErr.ErrPasswordChangeRequired
	}
	return actor.New(principal.ID, principal.Name, principal.Role)
}
