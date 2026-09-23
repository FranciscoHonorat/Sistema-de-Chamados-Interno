package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpapi "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/auth"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/backend/modules/tickets/internal/domain/domain-errors"
)

func newProtectedRouter(t *testing.T) *gin.Engine {
	t.Helper()
	admin, err := actor.New("admin-1", "Administradora", "admin")
	require.NoError(t, err)
	verifier := &outtest.TokenVerifier{Tokens: map[string]actor.Actor{"valid-token": admin}}

	r := gin.New()
	r.Use(httpapi.RequireAuthentication(auth.NewAuthenticateUseCase(verifier)))
	r.GET("/whoami", func(c *gin.Context) {
		c.String(http.StatusOK, httpapi.ActorFrom(c).ID())
	})
	return r
}

func TestRequireAuthentication(t *testing.T) {
	t.Run("should let a valid bearer token through and expose the actor", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		newProtectedRouter(t).ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "admin-1", w.Body.String())
	})

	for name, header := range map[string]string{
		"no authorization header": "",
		"a non bearer scheme":     "Basic valid-token",
		"an invalid bearer token": "Bearer forged",
		"an empty bearer token":   "Bearer ",
	} {
		t.Run("should answer 401 for "+name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
			if header != "" {
				req.Header.Set("Authorization", header)
			}
			w := httptest.NewRecorder()

			newProtectedRouter(t).ServeHTTP(w, req)

			require.Equal(t, http.StatusUnauthorized, w.Code)
			assert.JSONEq(t, `{"error":"unauthenticated"}`, w.Body.String())
		})
	}

	t.Run("should answer 403 while a temporary password was not replaced", func(t *testing.T) {
		verifier := &outtest.TokenVerifier{Err: domainErr.ErrPasswordChangeRequired}
		r := gin.New()
		r.Use(httpapi.RequireAuthentication(auth.NewAuthenticateUseCase(verifier)))
		r.GET("/whoami", func(c *gin.Context) { c.Status(http.StatusOK) })
		req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
		req.Header.Set("Authorization", "Bearer pending-token")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, `{"error":"password change required: set a new password before using the system"}`, w.Body.String())
	})
}
