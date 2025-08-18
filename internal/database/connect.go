package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectDB(connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return pool, nil
}
