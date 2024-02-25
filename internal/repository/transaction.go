package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TransactionRepos struct {
	pool *pgxpool.Pool
}

func NewTransactionsRepos(pool *pgxpool.Pool) Transaction {
	return &TransactionRepos{pool: pool}
}

func (r *TransactionRepos) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *TransactionRepos) Rollback(ctx context.Context, tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return tx.Rollback(ctx)
}
