// Package database contain low db instance wrappers and helpers
package database

import (
	"fmt"
	"time"
)

type Config struct {
	Host               string        `json:"host" yaml:"host" env:"HOST" validate:"notEmpty"`
	Port               int           `json:"port" yaml:"port" env:"PORT" default:"5432"`
	Database           string        `json:"database" yaml:"database" env:"DATABASE" validate:"notEmpty"`
	Username           string        `json:"username" yaml:"username" env:"USERNAME" validate:"notEmpty"`
	Password           string        `json:"password" yaml:"password" env:"PASSWORD" validate:"notEmpty"`
	MaxOpenConns       int           `json:"max_open_conns" yaml:"max_open_conns" env:"MAX_OPEN_CONNS" default:"10"`
	MaxConnLifetime    time.Duration `json:"max_conn_lifetime" yaml:"max_conn_lifetime" env:"MAX_CONN_LIFETIME" default:"1h"`
	SSL                bool          `json:"ssl" yaml:"ssl" env:"SSL" default:"false"`
	EnableLogger       bool          `json:"enable_logger" yaml:"enable_logger" env:"ENABLE_LOGGER" default:"true"`
	LoggerLevel        string        `json:"logger_level" yaml:"logger_level" env:"LOGGER_LEVEL" default:"debug"`
	StatementCacheMode string        `json:"statement_cache_mode" yaml:"statement_cache_mode" env:"STATEMENT_CACHE_MODE" default:"describe"`
}

func (c *Config) DSNFormat() string {
	sslmod := "disable"
	if c.SSL {
		sslmod = "enable"
	}

	statementCacheMode := "describe"
	if c.StatementCacheMode != "" {
		statementCacheMode = c.StatementCacheMode
	}

	return fmt.Sprintf(`host=%s port=%d dbname=%s user=%s password='%s' sslmode=%s statement_cache_mode=%s`,
		c.Host,
		c.Port,
		c.Database,
		c.Username,
		c.Password,
		sslmod,
		statementCacheMode)
}
