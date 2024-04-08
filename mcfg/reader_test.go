package mcfg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			type Cfg struct {
				Str string `env:"STR"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				key := "STR"
				value := "str_value"
				expected := "str_value"
				t.Setenv(key, value)
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Str)
			})
		})

		t.Run("default", func(t *testing.T) {
			type Cfg struct {
				Str string `env:"STR" default:"default_val"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				expected := "default_val"
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Str)
			})
		})
	})

	t.Run("slice string", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			type Cfg struct {
				SliceStr []string `env:"SLICE_STR"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				key := "SLICE_STR"
				value := "value1,value2"
				expected := []string{"value1", "value2"}
				t.Setenv(key, value)
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.SliceStr)
			})
		})

		t.Run("default", func(t *testing.T) {
			type Cfg struct {
				SliceStr []string `env:"SLICE_STR" default:"[val1,val2]"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				expected := []string{"val1", "val2"}
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.SliceStr)
			})
		})
	})

	t.Run("int", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			type Cfg struct {
				Int int `env:"INT"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				key := "INT"
				value := "10"
				expected := 10
				t.Setenv(key, value)
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Int)
			})
		})

		t.Run("default", func(t *testing.T) {
			type Cfg struct {
				Int int `env:"INT" default:"222"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				expected := 222
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Int)
			})
		})
	})

	t.Run("nested", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			type Cfg struct {
				Nested struct {
					Str string `env:"STR"`
				} `envPrefix:"NESTED_"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				key := "NESTED_STR"
				value := "nested_str_value"
				expected := "nested_str_value"
				t.Setenv(key, value)
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Nested.Str)
			})
		})

		t.Run("default", func(t *testing.T) {
			type Cfg struct {
				Nested struct {
					Str string `env:"STR" default:"default_val"`
				} `envPrefix:"NESTED_"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				expected := "default_val"
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Nested.Str)
			})
		})
	})

	t.Run("time.Duration", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			type Cfg struct {
				Dur time.Duration `env:"DUR"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				key := "DUR"
				value := "5m"
				expected := time.Minute * 5
				t.Setenv(key, value)
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Dur)
			})
		})

		t.Run("default", func(t *testing.T) {
			type Cfg struct {
				Dur time.Duration `env:"DUR" default:"7m"`
			}

			t.Run("env", func(t *testing.T) {
				var cfg Cfg
				expected := time.Minute * 7
				require.NoError(t, LoadFromEnv(&cfg))
				require.Equal(t, expected, cfg.Dur)
			})
		})
	})
}
