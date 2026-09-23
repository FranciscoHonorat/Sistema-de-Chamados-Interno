// Command sys-called is the whole Sys-Called backend in one binary.
//
//	sys-called [serve]    migrate (unless MIGRATE_ON_START=false) and serve
//	sys-called migrate    apply the pending migrations and exit
//	sys-called healthcheck [url]  exit 0 when the instance (local by default) is ready
//	sys-called version    print the version and commit
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/app"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/config"
	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/logging"
)

// Set at build time with -ldflags "-X main.version=... -X main.commit=...".
var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "healthcheck":
		os.Exit(healthcheck(os.Args[2:]))
	case "version":
		fmt.Println(version, commit)
		return
	}

	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	logger := logging.New(os.Stdout, cfg.LogLevel, cfg.LogFormat)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch command {
	case "serve":
		logger.Info("starting", "version", version, "commit", commit, "env", cfg.Env, "port", cfg.Port)
		err = app.Run(ctx, cfg, logger)
	case "migrate":
		err = app.Migrate(ctx, cfg, logger)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q (use serve, migrate or healthcheck)\n", command)
		os.Exit(2)
	}
	if err != nil {
		logger.Error("exiting with error", "command", command, "error", err)
		os.Exit(1)
	}
}

// healthcheck lets a distroless container (no shell, no curl) be probed by
// Docker's HEALTHCHECK, or check another instance when given its URL (the
// Helm test does that through the Service).
func healthcheck(args []string) int {
	url := "http://127.0.0.1:" + envOr("PORT", "8080") + "/readyz"
	if len(args) > 0 {
		url = args[0]
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url) //nolint:gosec,noctx // an operator-provided probe URL
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "not ready:", resp.Status)
		return 1
	}
	return 0
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
