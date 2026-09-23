package httpserver_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/httpserver"
)

func newServer(t *testing.T, staticDir string, checks map[string]httpserver.Check) http.Handler {
	t.Helper()
	s := httpserver.New(httpserver.Options{StaticDir: staticDir, Checks: checks, Logger: slog.New(slog.DiscardHandler)})
	s.API().GET("/things", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"things": strings.Repeat("x", 2048)}) })
	s.API().GET("/panic", func(*gin.Context) { panic("boom") })
	return s.Handler()
}

func staticDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "assets"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>spa</html>"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "assets", "app-123.js"), []byte("console.log(1)"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "favicon.svg"), []byte("<svg/>"), 0o600))
	return dir
}

func serve(h http.Handler, method, path string, headers ...string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestProbes(t *testing.T) {
	t.Run("liveness never looks at the dependencies", func(t *testing.T) {
		h := newServer(t, "", map[string]httpserver.Check{"database": func(context.Context) error { return errors.New("down") }})
		assert.Equal(t, http.StatusOK, serve(h, http.MethodGet, "/healthz").Code)
	})

	t.Run("readiness fails while a dependency is down, and says which", func(t *testing.T) {
		h := newServer(t, "", map[string]httpserver.Check{"database": func(context.Context) error { return errors.New("down") }})
		w := serve(h, http.MethodGet, "/readyz")
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), `"database":"down"`)
	})

	t.Run("readiness passes when every dependency answers", func(t *testing.T) {
		h := newServer(t, "", map[string]httpserver.Check{"database": func(context.Context) error { return nil }})
		assert.Equal(t, http.StatusOK, serve(h, http.MethodGet, "/readyz").Code)
	})
}

func TestSPA(t *testing.T) {
	h := newServer(t, staticDir(t), nil)

	t.Run("deep links get index.html, never cached", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/chamados/abc")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "<html>spa</html>", w.Body.String())
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	})

	t.Run("hashed assets are cached for good", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/assets/app-123.js")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "console.log(1)", w.Body.String())
		assert.Contains(t, w.Header().Get("Cache-Control"), "immutable")
	})

	t.Run("other files are served and revalidated", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/favicon.svg")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	})

	t.Run("path traversal stays inside the static directory", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/../../etc/passwd")
		assert.NotContains(t, w.Body.String(), "root:")
	})

	t.Run("unknown API paths are a JSON 404, not the SPA", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/api/unknown")
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"not found"}`, w.Body.String())
	})
}

func TestAPIOnlyWithoutStaticDir(t *testing.T) {
	h := newServer(t, "", nil)
	assert.Equal(t, http.StatusNotFound, serve(h, http.MethodGet, "/chamados").Code)
}

func TestMiddleware(t *testing.T) {
	h := newServer(t, "", nil)

	t.Run("sets the security headers", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/api/things")
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Contains(t, w.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'")
		assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
	})

	t.Run("sends HSTS behind a TLS-terminating proxy", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/api/things", "X-Forwarded-Proto", "https")
		assert.NotEmpty(t, w.Header().Get("Strict-Transport-Security"))
	})

	t.Run("propagates or creates a request id", func(t *testing.T) {
		assert.Equal(t, "abc", serve(h, http.MethodGet, "/api/things", "X-Request-ID", "abc").Header().Get("X-Request-ID"))
		assert.Len(t, serve(h, http.MethodGet, "/api/things").Header().Get("X-Request-ID"), 16)
	})

	t.Run("compresses responses for clients that accept gzip", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/api/things", "Accept-Encoding", "gzip")
		require.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
		reader, err := gzip.NewReader(w.Body)
		require.NoError(t, err)
		raw, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Contains(t, string(raw), `"things"`)
	})

	t.Run("turns a panic into a 500", func(t *testing.T) {
		w := serve(h, http.MethodGet, "/api/panic")
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())
	})
}
