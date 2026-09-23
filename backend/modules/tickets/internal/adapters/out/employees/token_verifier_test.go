package employees_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/contracts"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/out/employees"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
)

type fakeIdentity map[string]contracts.Principal

func (f fakeIdentity) VerifyAccessToken(_ context.Context, token string) (contracts.Principal, error) {
	principal, ok := f[token]
	if !ok {
		return contracts.Principal{}, assert.AnError
	}
	return principal, nil
}

func TestTokenVerifier(t *testing.T) {
	identity := fakeIdentity{
		"support-token": {ID: "agent-1", Name: "Ana Souza", Role: "support"},
		"weird-token":   {ID: "x-1", Name: "X", Role: "superuser"},
		"pending-token": {ID: "agent-2", Name: "Bruno Lima", Role: "support", MustChangePassword: true},
	}
	verifier := employees.NewTokenVerifier(identity)

	t.Run("should turn the principal into an actor", func(t *testing.T) {
		a, err := verifier.Verify(context.Background(), "support-token")

		require.NoError(t, err)
		assert.Equal(t, "agent-1", a.ID())
		assert.Equal(t, "Ana Souza", a.Name())
		assert.Equal(t, actor.RoleSupport, a.Role())
	})

	t.Run("should reject a token the employees module does not accept", func(t *testing.T) {
		_, err := verifier.Verify(context.Background(), "forged")

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should refuse someone who still has to replace a temporary password", func(t *testing.T) {
		_, err := verifier.Verify(context.Background(), "pending-token")

		assert.ErrorIs(t, err, domainErr.ErrPasswordChangeRequired)
	})

	t.Run("should reject a principal with a role the tickets module does not know", func(t *testing.T) {
		_, err := verifier.Verify(context.Background(), "weird-token")

		assert.Error(t, err)
	})
}
