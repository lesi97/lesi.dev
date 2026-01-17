package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type DB struct {
	*pgxpool.Pool
}


func Connect(logger *utils.Logger) (*DB, error) {
	url := os.Getenv("DATABASE_URL")
	
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Fatalf("Failed to parse pool config: %v", err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatalf("Failed to create connection pool: %v", err)
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func Disconnect(db *DB) {
	db.Close()
}
