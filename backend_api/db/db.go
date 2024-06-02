package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPool struct {
	*pgxpool.Pool
}

func Connect(ctx context.Context) (*DBPool, error) {
	pool, err := pgxpool.New(ctx, "postgres://mike:password@db:5432/db")
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	return &DBPool{Pool: pool}, nil
}

func (pool *DBPool) Disconnect() {
	pool.Close()
}
