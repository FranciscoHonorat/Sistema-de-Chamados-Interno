package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/backend/internal/platform/config"
)

func env(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := values[key]
		return v, ok
	}
}

func TestLoad(t *testing.T) {
	t.Run("should apply development defaults", func(t *testing.T) {
		cfg, err := config.Load(env(map[string]string{"DATABASE_URL": "postgres://db"}))

		require.NoError(t, err)
		assert.Equal(t, config.EnvDevelopment, cfg.Env)
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "9090", cfg.MetricsPort)
		assert.Equal(t, int32(10), cfg.DBMaxConns)
		assert.True(t, cfg.MigrateOnStart)
		assert.True(t, cfg.SeedDemoUsers)
		assert.Equal(t, "text", cfg.LogFormat)
		assert.Equal(t, 15*time.Second, cfg.ShutdownTimeout)
	})

	t.Run("should not seed demo users and log JSON in production", func(t *testing.T) {
		cfg, err := config.Load(env(map[string]string{
			"APP_ENV": "production", "DATABASE_URL": "postgres://db", "JWT_PRIVATE_KEY": "pem",
		}))

		require.NoError(t, err)
		assert.False(t, cfg.SeedDemoUsers)
		assert.Equal(t, "json", cfg.LogFormat)
	})

	t.Run("should require a signing key in production", func(t *testing.T) {
		_, err := config.Load(env(map[string]string{"APP_ENV": "production", "DATABASE_URL": "postgres://db"}))

		assert.ErrorContains(t, err, "JWT_PRIVATE_KEY is required")
	})

	t.Run("should report every invalid variable at once", func(t *testing.T) {
		_, err := config.Load(env(map[string]string{
			"APP_ENV": "staging", "MIGRATE_ON_START": "maybe", "DB_MAX_CONNS": "-1", "SHUTDOWN_TIMEOUT": "soon",
		}))

		require.Error(t, err)
		for _, fragment := range []string{"DATABASE_URL is required", "APP_ENV", "MIGRATE_ON_START", "DB_MAX_CONNS", "SHUTDOWN_TIMEOUT"} {
			assert.ErrorContains(t, err, fragment)
		}
	})
}
