package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

// TxManager wraps a pool and exposes WithTransaction for service-level use.
type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithTransaction starts a transaction, injects it into ctx, calls fn, and
// commits on success or rolls back on any error.
func (m *TxManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err := fn(context.WithValue(ctx, txKey{}, DBTX(tx))); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ExtractDB returns the transaction stored in ctx by WithTransaction, or falls
// back to the raw pool when called outside a transactional context.
func ExtractDB(ctx context.Context, pool *pgxpool.Pool) DBTX {
	if tx, ok := ctx.Value(txKey{}).(DBTX); ok {
		return tx
	}
	return pool
}

// InjectTx returns a context with tx embedded so that ExtractDB will use it.
// This is intentionally exported for use in test helpers only.
func InjectTx(ctx context.Context, tx DBTX) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
