package mlog

import (
	"testing"

	"github.com/krivenkov/pkg/mcfg"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("env", func(t *testing.T) {
		var (
			cfg    Config
			level  = "WARN"
			encode = "console"
			caller = "full"
		)
		t.Setenv("LEVEL", level)
		t.Setenv("ENCODE", encode)
		t.Setenv("CALLER", caller)
		err := mcfg.LoadFromEnv(&cfg)
		require.NoError(t, err)
		require.Equal(t, Config{
			Level:  level,
			Encode: encode,
			Caller: caller,
		}, cfg)
	})

	t.Run("defaults", func(t *testing.T) {
		var cfg Config
		err := mcfg.LoadFromEnv(&cfg)
		require.NoError(t, err)
		require.Equal(t, Config{
			Level:  "INFO",
			Encode: "json",
		}, cfg)
	})
}
