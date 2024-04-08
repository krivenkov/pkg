package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/krivenkov/pkg/global"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewPgxPool(cfg Config, info global.Info, logger *zap.Logger, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	logger = logger.With(zap.String("module", "pgx"))

	dbPoolCfg, err := pgxpool.ParseConfig(cfg.DSNFormat())
	if err != nil {
		return nil, err
	}

	dbPoolCfg.ConnConfig.RuntimeParams["application_name"] = info.AppName

	if cfg.EnableLogger {
		logLevel, errLog := pgx.LogLevelFromString(cfg.LoggerLevel)
		if errLog != nil {
			return nil, fmt.Errorf("parse LogLevel: %w", errLog)
		}
		dbPoolCfg.ConnConfig.LogLevel = logLevel
		dbPoolCfg.ConnConfig.Logger = zapadapter.NewLogger(logger)
	}

	if cfg.MaxConnLifetime > 0 {
		dbPoolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	}

	if cfg.MaxOpenConns > 0 {
		dbPoolCfg.MaxConns = int32(cfg.MaxOpenConns)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), dbPoolCfg)
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return pool.Ping(ctx)
			},
			OnStop: func(_ context.Context) error {
				pool.Close()
				return nil
			},
		},
	)

	return pool, nil
}
