// Package dbtest gives the Postgres integration tests a pool on a migrated
// schema. The tests are skipped unless SYS_CALLED_TEST_DATABASE_URL points to a
// database they may write to (CI runs them against a service container).
package dbtest

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/database"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/migrate"
)

const EnvVar = "SYS_CALLED_TEST_DATABASE_URL"

// DSN returns the test database URL, skipping the test when it is not set.
func DSN(t testing.TB) string {
	t.Helper()
	dsn := os.Getenv(EnvVar)
	if dsn == "" {
		t.Skip(EnvVar + " not set; skipping Postgres integration tests")
	}
	return dsn
}

// Pool opens a pool pinned to schema, after applying migrations to it.
func Pool(t testing.TB, schema string, migrations fs.FS) *pgxpool.Pool {
	t.Helper()

	logger := slog.New(slog.DiscardHandler)
	ctx := context.Background()

	pool, err := database.Open(ctx, DSN(t), database.Options{Schema: schema, ConnectTimeout: 10 * time.Second}, logger)
	if err != nil {
		t.Fatalf("opening test database: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := migrate.Run(ctx, pool, migrate.Source{Schema: schema, Files: migrations}, logger); err != nil {
		t.Fatalf("migrating test database: %v", err)
	}
	return pool
}
