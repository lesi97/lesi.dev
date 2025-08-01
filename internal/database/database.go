package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Supabase = pgxpool.Pool

func Connect(logger *log.Logger) (*Supabase, error) {
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

	return pool, nil
}

func Disconnect(supabase *Supabase) {
	supabase.Close()
}
