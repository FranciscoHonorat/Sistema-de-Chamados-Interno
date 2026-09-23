// Package config reads the app configuration from the environment, validates
// it once at startup and fails with every problem at the same time.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type Config struct {
	Env             string
	Port            string
	MetricsPort     string
	DatabaseURL     string
	DBMaxConns      int32
	JWTPrivateKey   string
	StaticDir       string
	MigrateOnStart  bool
	SeedDemoUsers   bool
	CacheTickets    bool
	LogLevel        string
	LogFormat       string
	ShutdownTimeout time.Duration
}

func (c Config) IsProduction() bool {
	return c.Env == EnvProduction
}

// Load builds the configuration from lookup (os.LookupEnv in production).
func Load(lookup func(string) (string, bool)) (Config, error) {
	r := reader{lookup: lookup}

	cfg := Config{
		Env:             r.oneOf("APP_ENV", EnvDevelopment, EnvDevelopment, EnvProduction),
		Port:            r.string("PORT", "8080"),
		MetricsPort:     r.string("METRICS_PORT", "9090"),
		DatabaseURL:     r.required("DATABASE_URL"),
		JWTPrivateKey:   r.string("JWT_PRIVATE_KEY", ""),
		StaticDir:       r.string("STATIC_DIR", ""),
		MigrateOnStart:  r.bool("MIGRATE_ON_START", true),
		CacheTickets:    r.bool("TICKETS_CACHE_ENABLED", true),
		LogLevel:        r.oneOf("LOG_LEVEL", "info", "debug", "info", "warn", "error"),
		ShutdownTimeout: r.duration("SHUTDOWN_TIMEOUT", 15*time.Second),
	}
	cfg.DBMaxConns = r.int32("DB_MAX_CONNS", 10)
	cfg.LogFormat = r.oneOf("LOG_FORMAT", defaultLogFormat(cfg.Env), "text", "json")
	// The demo users all share a known password: never create them in
	// production unless explicitly asked to.
	cfg.SeedDemoUsers = r.bool("SEED_DEMO_USERS", !cfg.IsProduction())

	if cfg.IsProduction() && cfg.JWTPrivateKey == "" {
		r.errs = append(r.errs, errors.New("JWT_PRIVATE_KEY is required when APP_ENV=production"))
	}

	if len(r.errs) > 0 {
		return Config{}, fmt.Errorf("invalid configuration: %w", errors.Join(r.errs...))
	}
	return cfg, nil
}

func FromEnv() (Config, error) {
	return Load(os.LookupEnv)
}

func defaultLogFormat(env string) string {
	if env == EnvProduction {
		return "json"
	}
	return "text"
}

type reader struct {
	lookup func(string) (string, bool)
	errs   []error
}

func (r *reader) string(key, fallback string) string {
	if v, ok := r.lookup(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return fallback
}

func (r *reader) required(key string) string {
	v := r.string(key, "")
	if v == "" {
		r.errs = append(r.errs, fmt.Errorf("%s is required", key))
	}
	return v
}

func (r *reader) oneOf(key, fallback string, allowed ...string) string {
	v := strings.ToLower(r.string(key, fallback))
	for _, a := range allowed {
		if v == a {
			return v
		}
	}
	r.errs = append(r.errs, fmt.Errorf("%s must be one of %s, got %q", key, strings.Join(allowed, ", "), v))
	return fallback
}

func (r *reader) bool(key string, fallback bool) bool {
	raw := r.string(key, "")
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		r.errs = append(r.errs, fmt.Errorf("%s must be a boolean, got %q", key, raw))
		return fallback
	}
	return v
}

func (r *reader) int32(key string, fallback int32) int32 {
	raw := r.string(key, "")
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || v <= 0 {
		r.errs = append(r.errs, fmt.Errorf("%s must be a positive integer, got %q", key, raw))
		return fallback
	}
	return int32(v)
}

func (r *reader) duration(key string, fallback time.Duration) time.Duration {
	raw := r.string(key, "")
	if raw == "" {
		return fallback
	}
	v, err := time.ParseDuration(raw)
	if err != nil || v <= 0 {
		r.errs = append(r.errs, fmt.Errorf("%s must be a positive duration like 15s, got %q", key, raw))
		return fallback
	}
	return v
}
