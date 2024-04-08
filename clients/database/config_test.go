package database

import (
	"testing"
	"time"

	"github.com/krivenkov/pkg/mcfg"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("env", func(t *testing.T) {
		t.Setenv("HOST", "host.pg")
		t.Setenv("PORT", "3333")
		t.Setenv("DATABASE", "dbname")
		t.Setenv("USERNAME", "username")
		t.Setenv("PASSWORD", "password")
		t.Setenv("MAX_OPEN_CONNS", "33")
		t.Setenv("MAX_CONN_LIFETIME", "33m")
		t.Setenv("SSL", "true")
		t.Setenv("ENABLE_LOGGER", "false")
		t.Setenv("LOGGER_LEVEL", "warn")
		t.Setenv("STATEMENT_CACHE_MODE", "prepare")

		var cfg Config
		err := mcfg.LoadFromEnv(&cfg)
		require.NoError(t, err)
		require.Equal(t, Config{
			Host:               "host.pg",
			Port:               3333,
			Database:           "dbname",
			Username:           "username",
			Password:           "password",
			MaxOpenConns:       33,
			MaxConnLifetime:    time.Minute * 33,
			SSL:                true,
			EnableLogger:       false,
			LoggerLevel:        "warn",
			StatementCacheMode: "prepare",
		}, cfg)
	})

	t.Run("defaults", func(t *testing.T) {
		t.Setenv("HOST", "host.pg")
		t.Setenv("DATABASE", "dbname")
		t.Setenv("USERNAME", "username")
		t.Setenv("PASSWORD", "password")

		var cfg Config
		err := mcfg.LoadFromEnv(&cfg)
		require.NoError(t, err)
		require.Equal(t, Config{
			Host:               "host.pg",
			Port:               5432,
			Database:           "dbname",
			Username:           "username",
			Password:           "password",
			MaxOpenConns:       10,
			MaxConnLifetime:    time.Hour,
			SSL:                false,
			EnableLogger:       true,
			LoggerLevel:        "debug",
			StatementCacheMode: "describe",
		}, cfg)
	})
}
