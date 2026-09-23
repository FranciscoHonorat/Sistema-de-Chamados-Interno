// Package database opens the connection pools of the monolith. There is a
// single Postgres database, and each module owns one schema in it: its pool
// has search_path pinned to that schema, so a module's SQL cannot reach
// another module's tables by accident.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var validSchema = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

type Options struct {
	Schema   string
	MaxConns int32
	// ConnectTimeout bounds how long Open keeps retrying while the database
	// is not accepting connections yet (a container that is still starting).
	ConnectTimeout time.Duration
}

func Open(ctx context.Context, dsn string, opts Options, logger *slog.Logger) (*pgxpool.Pool, error) {
	if !validSchema.MatchString(opts.Schema) {
		return nil, fmt.Errorf("database: invalid schema name %q", opts.Schema)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("database: invalid DATABASE_URL: %w", err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = opts.Schema
	config.ConnConfig.RuntimeParams["application_name"] = "sys-called/" + opts.Schema
	if opts.MaxConns > 0 {
		config.MaxConns = opts.MaxConns
	}
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	timeout := opts.ConnectTimeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	if err := waitUntilReachable(ctx, pool, timeout, logger.With("schema", opts.Schema)); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func waitUntilReachable(ctx context.Context, pool *pgxpool.Pool, timeout time.Duration, logger *slog.Logger) error {
	deadline := time.Now().Add(timeout)
	backoff := 250 * time.Millisecond
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := pool.Ping(pingCtx)
		cancel()
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("database: not reachable after %s: %w", timeout, err)
		}
		logger.Warn("database not reachable yet, retrying", "error", err, "retry_in", backoff)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, 5*time.Second)
	}
}
