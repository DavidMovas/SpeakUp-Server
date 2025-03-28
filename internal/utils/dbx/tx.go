package dbx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey struct{}

var TxContextKey contextKey = struct{}{}

func WithTransaction(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxContextKey, tx)
}

func FromContext(ctx context.Context, def Queryable) Queryable {
	if tx, ok := ctx.Value(TxContextKey).(pgx.Tx); ok {
		return tx
	}
	return def
}

func InTransaction(ctx context.Context, db *pgxpool.Pool, fn func(ctx2 context.Context, tx pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctx = WithTransaction(ctx, tx)
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				err = errors.Join(err, fmt.Errorf("failed to rollback transaction: %w", txErr))
			}
		} else {
			cerr := tx.Commit(ctx)
			if cerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", cerr)
			}
		}
	}()

	return fn(ctx, tx)
}
