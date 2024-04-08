package database

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/krivenkov/pkg/txer"
	pgx "pkg/mod/github.com/jackc/pgx/v4"
)

var ErrTxAlready = errors.New("tx is already running in the context")

type TXer struct {
	db *pgxpool.Pool
}

func NewTXer(db *pgxpool.Pool) *TXer {
	return &TXer{
		db: db,
	}
}

func NewTXerFX(db *pgxpool.Pool) (*TXer, txer.TXer) {
	v := NewTXer(db)
	return v, v
}

type ctxKey string

const (
	ctxKeyTx ctxKey = "TransactionContextKey"
)

func (t *TXer) WithTX(ctx context.Context, f func(context.Context) error) error {
	if tx := txFromCtx(ctx); tx != nil {
		return ErrTxAlready
	}

	return t.db.BeginFunc(ctx, func(tx pgx.Tx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("recover: %v stack: %v", r, string(debug.Stack()))
			}
		}()

		return f(ctxWithTx(ctx, tx))
	})
}

func (t *TXer) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	if tx := txFromCtx(ctx); tx != nil {
		return f(tx)
	}

	return t.db.BeginFunc(ctx, f)
}

func ctxWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx, tx)
}

func txFromCtx(ctx context.Context) pgx.Tx {
	if tx, exists := ctx.Value(ctxKeyTx).(pgx.Tx); exists {
		return tx
	}

	return nil
}
