// Package httpserver is the single HTTP entry point of the monolith: it serves
// the modules' APIs under /api, the built frontend, the health probes and,
// on a separate port, the Prometheus metrics.
package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Options struct {
	Port        string
	MetricsPort string
	StaticDir   string
	Checks      map[string]Check
	Logger      *slog.Logger
}

type Server struct {
	engine   *gin.Engine
	api      *gin.RouterGroup
	registry *prometheus.Registry
	opts     Options
}

func New(opts Options) *Server {
	gin.SetMode(gin.ReleaseMode)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	engine := gin.New()
	engine.HandleMethodNotAllowed = false
	// The app sits behind a single reverse proxy (the ingress) at most; trust
	// X-Forwarded-For only from private networks.
	_ = engine.SetTrustedProxies([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1/32", "::1/128"})
	engine.Use(
		requestID(),
		recovery(opts.Logger),
		accessLog(opts.Logger),
		newHTTPMetrics(registry).middleware(),
		securityHeaders(),
		limitBody(),
	)
	registerHealth(engine, opts.Checks)

	if staticDirUsable(opts.StaticDir) {
		engine.NoRoute(serveSPA(opts.StaticDir))
	} else {
		if opts.StaticDir != "" {
			opts.Logger.Warn("STATIC_DIR has no index.html, serving the API only", "static_dir", opts.StaticDir)
		}
		engine.NoRoute(func(c *gin.Context) { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}) })
	}

	return &Server{engine: engine, api: engine.Group("/api"), registry: registry, opts: opts}
}

// API is the /api group every module mounts its routes under.
func (s *Server) API() *gin.RouterGroup {
	return s.api
}

// Registry lets the modules register their own metrics.
func (s *Server) Registry() prometheus.Registerer {
	return s.registry
}

// Handler is the full HTTP handler of the app, with response compression.
func (s *Server) Handler() http.Handler {
	return gzhttp.GzipHandler(s.engine)
}

// Run serves until ctx is cancelled, then drains the in-flight requests for up
// to shutdownTimeout.
func (s *Server) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	app := &http.Server{
		Addr:              ":" + s.opts.Port,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ErrorLog:          slog.NewLogLogger(s.opts.Logger.Handler(), slog.LevelWarn),
	}
	metrics := &http.Server{
		Addr:              ":" + s.opts.MetricsPort,
		Handler:           promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errs := make(chan error, 2)
	for _, srv := range []*http.Server{app, metrics} {
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errs <- err
			}
		}()
	}
	s.opts.Logger.Info("http server listening", "port", s.opts.Port, "metrics_port", s.opts.MetricsPort)

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
	}

	s.opts.Logger.Info("shutting down http server", "timeout", shutdownTimeout)
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()
	return errors.Join(app.Shutdown(shutdownCtx), metrics.Shutdown(shutdownCtx))
}
