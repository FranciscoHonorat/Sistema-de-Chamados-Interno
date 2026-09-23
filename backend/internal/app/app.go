// Package app is the composition root of the monolith: it opens the database,
// migrates each module's schema, builds the modules, wires them together through
// their public contracts and the event bus, and runs the HTTP server and the
// background workers until the process is asked to stop.
package app

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/config"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/database"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/eventbus"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/httpserver"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/migrate"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/employees"
	"github.com/franciscoHonorat/Sys-Called/backend/modules/tickets"
)

// module describes what the composition root needs to know about a module
// before building it: the schema it owns and the migrations of that schema.
type module struct {
	schema     string
	migrations fs.FS
}

var modules = []module{
	{schema: employees.Schema, migrations: employees.Migrations()},
	{schema: tickets.Schema, migrations: tickets.Migrations()},
}

type pools map[string]*pgxpool.Pool

func (p pools) Close() {
	for _, pool := range p {
		pool.Close()
	}
}

func openPools(ctx context.Context, cfg config.Config, logger *slog.Logger) (pools, error) {
	opened := pools{}
	for _, m := range modules {
		pool, err := database.Open(ctx, cfg.DatabaseURL, database.Options{Schema: m.schema, MaxConns: cfg.DBMaxConns}, logger)
		if err != nil {
			opened.Close()
			return nil, err
		}
		opened[m.schema] = pool
	}
	return opened, nil
}

func migrateAll(ctx context.Context, p pools, logger *slog.Logger) error {
	for _, m := range modules {
		if err := migrate.Run(ctx, p[m.schema], migrate.Source{Schema: m.schema, Files: m.migrations}, logger); err != nil {
			return err
		}
	}
	logger.Info("database schema up to date")
	return nil
}

// Migrate applies every module's pending migrations and returns.
func Migrate(ctx context.Context, cfg config.Config, logger *slog.Logger) error {
	p, err := openPools(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer p.Close()
	return migrateAll(ctx, p, logger)
}

// App is the assembled monolith: the modules wired together and mounted on
// the HTTP server.
type App struct {
	cfg       config.Config
	logger    *slog.Logger
	pools     pools
	server    *httpserver.Server
	employees *employees.Module
}

// New opens the database, migrates it when configured to, and assembles the
// modules. Close releases what it opened.
func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	p, err := openPools(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}
	a, err := assemble(ctx, cfg, logger, p)
	if err != nil {
		p.Close()
		return nil, err
	}
	return a, nil
}

func assemble(ctx context.Context, cfg config.Config, logger *slog.Logger, p pools) (*App, error) {
	if cfg.MigrateOnStart {
		if err := migrateAll(ctx, p, logger); err != nil {
			return nil, err
		}
	}

	bus := eventbus.New()

	employeesModule, err := employees.New(ctx, employees.Config{
		Pool:              p[employees.Schema],
		Bus:               bus,
		JWTPrivateKey:     cfg.JWTPrivateKey,
		RequireSigningKey: cfg.IsProduction(),
		SeedDemoUsers:     cfg.SeedDemoUsers,
		CookiePath:        "/api/employees/auth",
		Logger:            logger.With("module", "employees"),
	})
	if err != nil {
		return nil, err
	}
	ticketsModule := tickets.New(tickets.Config{
		Pool:         p[tickets.Schema],
		Bus:          bus,
		Identity:     employeesModule,
		CacheTickets: cfg.CacheTickets,
		Logger:       logger.With("module", "tickets"),
	})

	server := httpserver.New(httpserver.Options{
		Port:        cfg.Port,
		MetricsPort: cfg.MetricsPort,
		StaticDir:   cfg.StaticDir,
		Checks: map[string]httpserver.Check{
			"database": func(ctx context.Context) error { return p[employees.Schema].Ping(ctx) },
		},
		Logger: logger,
	})
	employeesModule.RegisterRoutes(server.API().Group("/employees"))
	ticketsModule.RegisterRoutes(server.API().Group("/tickets"))

	return &App{cfg: cfg, logger: logger, pools: p, server: server, employees: employeesModule}, nil
}

// Handler is the app's HTTP handler, for tests that serve it with httptest.
func (a *App) Handler() http.Handler {
	return a.server.Handler()
}

// StartWorkers runs the modules' background work until ctx ends; the returned
// function waits for it to finish.
func (a *App) StartWorkers(ctx context.Context) (wait func()) {
	var workers sync.WaitGroup
	workers.Go(func() { a.employees.Run(ctx) })
	return workers.Wait
}

// Serve runs the HTTP server and the workers until ctx is cancelled, then
// drains the requests in flight and stops the workers.
func (a *App) Serve(ctx context.Context) error {
	workersCtx, stopWorkers := context.WithCancel(ctx)
	waitWorkers := a.StartWorkers(workersCtx)

	err := a.server.Run(ctx, a.cfg.ShutdownTimeout)
	stopWorkers()
	waitWorkers()

	if err != nil {
		return fmt.Errorf("http server: %w", err)
	}
	a.logger.Info("stopped cleanly")
	return nil
}

func (a *App) Close() {
	a.pools.Close()
}

// Run assembles the app and serves it until ctx is cancelled.
func Run(ctx context.Context, cfg config.Config, logger *slog.Logger) error {
	a, err := New(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer a.Close()
	return a.Serve(ctx)
}
