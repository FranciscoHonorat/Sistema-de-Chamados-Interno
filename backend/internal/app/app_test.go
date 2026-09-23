package app_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/app"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/config"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/database/dbtest"
)

// TestMonolith boots the whole app against a real Postgres and drives it over
// HTTP: it proves the modules are wired together, that the employees module's
// events reach the tickets module through the outbox and the bus, and that the
// frontend is served next to the API.
func TestMonolith(t *testing.T) {
	dsn := dbtest.DSN(t)

	static := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(static, "index.html"), []byte(`<div id="app"></div>`), 0o600))

	cfg, err := config.Load(func(key string) (string, bool) {
		v, ok := map[string]string{"DATABASE_URL": dsn, "STATIC_DIR": static}[key]
		return v, ok
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	logger := slog.New(slog.DiscardHandler)

	// Migrating again must be a no-op.
	require.NoError(t, app.Migrate(ctx, cfg, logger))

	a, err := app.New(ctx, cfg, logger)
	require.NoError(t, err)
	t.Cleanup(a.Close)
	wait := a.StartWorkers(ctx)
	t.Cleanup(wait)
	t.Cleanup(cancel)

	srv := httptest.NewServer(a.Handler())
	t.Cleanup(srv.Close)

	t.Run("probes", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, get(t, srv.URL+"/healthz", "").StatusCode)
		assert.Equal(t, http.StatusOK, get(t, srv.URL+"/readyz", "").StatusCode)
	})

	t.Run("the frontend is served for any non-API path", func(t *testing.T) {
		resp := get(t, srv.URL+"/chamados/123", "")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, body(t, resp), `<div id="app">`)
		assert.NotEmpty(t, resp.Header.Get("Content-Security-Policy"))

		assert.Equal(t, http.StatusNotFound, get(t, srv.URL+"/api/nothing-here", "").StatusCode)
	})

	t.Run("the API refuses anonymous calls", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, get(t, srv.URL+"/api/tickets/tickets", "").StatusCode)
	})

	login := post(t, srv.URL+"/api/employees/auth/login", "", `{"username":"admin","password":"senha123"}`)
	require.Equal(t, http.StatusOK, login.StatusCode)
	var session struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.Unmarshal([]byte(body(t, login)), &session))
	token := session.AccessToken

	t.Run("the refresh cookie is scoped to the auth routes", func(t *testing.T) {
		var cookie *http.Cookie
		for _, c := range login.Cookies() {
			if c.Name == "refresh_token" {
				cookie = c
			}
		}
		require.NotNil(t, cookie)
		assert.Equal(t, "/api/employees/auth", cookie.Path)
		assert.True(t, cookie.HttpOnly)
	})

	t.Run("support agents reach the tickets module through the event bus", func(t *testing.T) {
		require.Eventually(t, func() bool {
			var agents []map[string]any
			_ = json.Unmarshal([]byte(body(t, get(t, srv.URL+"/api/tickets/responsibles", token))), &agents)
			return len(agents) >= 3
		}, 10*time.Second, 100*time.Millisecond)
	})

	t.Run("a ticket goes through opening and automatic assignment", func(t *testing.T) {
		opened := post(t, srv.URL+"/api/tickets/tickets", token, `{"title":"Monitor","description":"Não liga","priority":"High"}`)
		require.Equal(t, http.StatusCreated, opened.StatusCode, body(t, opened))
		var ticket struct {
			ID string `json:"ticket_id"`
		}
		require.NoError(t, json.Unmarshal([]byte(body(t, opened)), &ticket))

		assigned := post(t, srv.URL+"/api/tickets/tickets/"+ticket.ID+"/assign/auto", token, "")
		require.Equal(t, http.StatusOK, assigned.StatusCode)
		assert.Contains(t, body(t, assigned), `"assignee_id"`)

		assert.Contains(t, body(t, get(t, srv.URL+"/api/tickets/tickets", token)), ticket.ID)
	})

	t.Run("a sign-up notifies the administrators through the event bus", func(t *testing.T) {
		username := "e2e-" + time.Now().Format("150405.000000")
		signup := post(t, srv.URL+"/api/employees/auth/signup", "", `{"name":"Pessoa E2E","username":"`+username+`","password":"senha-forte"}`)
		require.Equal(t, http.StatusCreated, signup.StatusCode, body(t, signup))

		require.Eventually(t, func() bool {
			return strings.Contains(body(t, get(t, srv.URL+"/api/tickets/notifications", token)), "Pessoa E2E")
		}, 10*time.Second, 100*time.Millisecond)
	})
}

func get(t *testing.T, url, token string) *http.Response {
	t.Helper()
	return do(t, http.MethodGet, url, token, "")
}

func post(t *testing.T, url, token, payload string) *http.Response {
	t.Helper()
	return do(t, http.MethodPost, url, token, payload)
}

func do(t *testing.T, method, url, token, payload string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	require.NoError(t, err)
	if payload != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

// body reads the response body once and keeps it for later calls.
func body(t *testing.T, resp *http.Response) string {
	t.Helper()
	if cached, ok := bodies[resp]; ok {
		return cached
	}
	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodies[resp] = string(raw)
	return string(raw)
}

var bodies = map[*http.Response]string{}
