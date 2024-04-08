package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/krivenkov/pkg/mlog"
	"go.uber.org/zap"
)

func AcquireTX(ctx context.Context, pool *pgxpool.Pool, f func(tx pgx.Tx) error) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("db acquire: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	if err = f(tx); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			mlog.FromContext(ctx).Error("failed rollback tx", zap.Error(errRollback))
		}

		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed commit: %w", err)
	}

	return nil
}
