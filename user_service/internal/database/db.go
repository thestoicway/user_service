package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	DbPool *pgxpool.Pool
}

func NewDB(postgresUrl string) (*DB, error) {
	dbpool, err := pgxpool.New(context.Background(), postgresUrl)

	if err != nil {
		return nil, err
	}

	return &DB{
		DbPool: dbpool,
	}, nil
}
