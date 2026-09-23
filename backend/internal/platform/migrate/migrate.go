// Package migrate applies the SQL migrations each module embeds in its binary.
// Every module migrates its own schema and keeps its own history table, and a
// Postgres advisory lock makes concurrent runs (several replicas starting at
// once, or the migrate command racing a server) apply each file exactly once.
package migrate

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"path"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// lockID is an arbitrary constant shared by every instance of the app.
const lockID int64 = 0x5359_5343_414c_4c // "SYSCALL"

type Source struct {
	Schema string
	Files  fs.FS
}

func Run(ctx context.Context, pool *pgxpool.Pool, src Source, logger *slog.Logger) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, lockID); err != nil {
		return fmt.Errorf("migrate: acquiring lock: %w", err)
	}
	defer func() {
		_, _ = conn.Exec(context.WithoutCancel(ctx), `SELECT pg_advisory_unlock($1)`, lockID)
	}()

	schema := pgx.Identifier{src.Schema}.Sanitize()
	if _, err := conn.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS `+schema); err != nil {
		return fmt.Errorf("migrate: creating schema %s: %w", src.Schema, err)
	}
	if _, err := conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS `+schema+`.schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`); err != nil {
		return fmt.Errorf("migrate: creating history table: %w", err)
	}

	files, err := sqlFiles(src.Files)
	if err != nil {
		return err
	}

	for _, name := range files {
		applied, err := isApplied(ctx, conn.Conn(), schema, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		body, err := fs.ReadFile(src.Files, name)
		if err != nil {
			return err
		}
		if err := apply(ctx, conn.Conn(), schema, name, string(body)); err != nil {
			return fmt.Errorf("migrate: %s/%s: %w", src.Schema, name, err)
		}
		logger.Info("migration applied", "schema", src.Schema, "version", name)
	}
	return nil
}

func sqlFiles(files fs.FS) ([]string, error) {
	var names []string
	err := fs.WalkDir(files, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(p, ".sql") {
			names = append(names, p)
		}
		return nil
	})
	sort.Slice(names, func(i, j int) bool { return path.Base(names[i]) < path.Base(names[j]) })
	return names, err
}

func isApplied(ctx context.Context, conn *pgx.Conn, schema, name string) (bool, error) {
	var exists bool
	err := conn.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM `+schema+`.schema_migrations WHERE version = $1)`, path.Base(name)).Scan(&exists)
	return exists, err
}

func apply(ctx context.Context, conn *pgx.Conn, schema, name, body string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `SET LOCAL search_path TO `+schema); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, body); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO `+schema+`.schema_migrations (version) VALUES ($1)`, path.Base(name)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
