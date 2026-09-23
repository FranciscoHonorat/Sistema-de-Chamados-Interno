package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/employees/internal/domain/domain-errors"
)

func loggedInAdmin(t *testing.T) (*outtest.EmployeeRepository, *application.Sessions, *outtest.RefreshTokenStore, string) {
	t.Helper()
	repo := repoWithAdmin()
	store := outtest.NewRefreshTokenStore()
	sessions := newTestSessions(store)

	login, err := application.NewLoginUseCase(repo, &outtest.PasswordHasher{}, sessions).
		Execute(context.Background(), application.LoginInput{Username: "admin", Password: "secret"})
	require.NoError(t, err)

	return repo, sessions, store, login.RefreshToken
}

func TestRefreshSessionUseCase(t *testing.T) {
	t.Run("should rotate a valid refresh token into a brand new session", func(t *testing.T) {
		repo, sessions, store, refreshToken := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, sessions)

		output, err := uc.Execute(context.Background(), refreshToken)

		require.NoError(t, err)
		assert.Equal(t, "token-for-admin-1", output.AccessToken)
		assert.Equal(t, "refresh-2", output.RefreshToken)
		assert.True(t, store.Tokens["sha:"+refreshToken].IsRevoked())
		assert.False(t, store.Tokens["sha:refresh-2"].IsRevoked())
	})

	t.Run("should reject an unknown refresh token", func(t *testing.T) {
		repo, sessions, store, _ := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, sessions)

		_, err := uc.Execute(context.Background(), "forged-token")

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Len(t, store.Tokens, 1)
	})

	t.Run("should let a second tab reuse a token that rotated seconds ago", func(t *testing.T) {
		repo, sessions, store, refreshToken := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, sessions)
		first, err := uc.Execute(context.Background(), refreshToken)
		require.NoError(t, err)

		late := application.NewRefreshSessionUseCase(repo, newTestSessionsAt(store, testNow.Add(application.RefreshReuseWindow)))
		second, err := late.Execute(context.Background(), refreshToken)

		require.NoError(t, err)
		assert.NotEqual(t, first.RefreshToken, second.RefreshToken)
		assert.False(t, store.Tokens["sha:"+first.RefreshToken].IsRevoked())
		assert.False(t, store.Tokens["sha:"+second.RefreshToken].IsRevoked())
	})

	t.Run("should revoke every session of the employee when a rotated token is replayed after the reuse window", func(t *testing.T) {
		repo, sessions, store, stolenToken := loggedInAdmin(t)
		_, err := application.NewRefreshSessionUseCase(repo, sessions).Execute(context.Background(), stolenToken)
		require.NoError(t, err)

		replay := application.NewRefreshSessionUseCase(repo, newTestSessionsAt(store, testNow.Add(application.RefreshReuseWindow+time.Second)))
		_, err = replay.Execute(context.Background(), stolenToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		for hash, token := range store.Tokens {
			assert.True(t, token.IsRevoked(), hash)
		}
	})

	t.Run("should treat a logged out token as a replay even right after the logout", func(t *testing.T) {
		repo, sessions, store, refreshToken := loggedInAdmin(t)
		require.NoError(t, sessions.End(context.Background(), refreshToken))

		output, err := application.NewRefreshSessionUseCase(repo, sessions).Execute(context.Background(), refreshToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Empty(t, output.AccessToken)
		assert.Len(t, store.Tokens, 1)
	})

	t.Run("should close the reuse window once every session is revoked", func(t *testing.T) {
		repo, sessions, store, refreshToken := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, sessions)
		_, err := uc.Execute(context.Background(), refreshToken)
		require.NoError(t, err)
		require.NoError(t, sessions.EndAll(context.Background(), "admin-1"))

		_, err = uc.Execute(context.Background(), refreshToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		for hash, token := range store.Tokens {
			assert.True(t, token.IsRevoked(), hash)
		}
	})

	t.Run("should reject the refresh token of an employee that no longer exists", func(t *testing.T) {
		repo, sessions, _, refreshToken := loggedInAdmin(t)
		repo.Employees = nil
		uc := application.NewRefreshSessionUseCase(repo, sessions)

		output, err := uc.Execute(context.Background(), refreshToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Empty(t, output.AccessToken)
	})

	t.Run("should reject an expired refresh token", func(t *testing.T) {
		repo, _, store, refreshToken := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, newTestSessionsAt(store, testNow.Add(testRefreshTTL)))

		_, err := uc.Execute(context.Background(), refreshToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Len(t, store.Tokens, 1)
	})
}
