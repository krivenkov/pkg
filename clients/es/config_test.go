package es

import (
	"testing"

	"github.com/krivenkov/pkg/mcfg"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("env", func(t *testing.T) {
		t.Setenv("ADDRESSES", "host1.es:9200,host2.es:9200")
		t.Setenv("USERNAME", "username")
		t.Setenv("PASSWORD", "password")

		var cfg Config
		err := mcfg.LoadFromEnv(&cfg)
		require.NoError(t, err)
		require.Equal(t, Config{
			Addresses: []string{"host1.es:9200", "host2.es:9200"},
			Username:  "username",
			Password:  "password",
		}, cfg)
	})
}
