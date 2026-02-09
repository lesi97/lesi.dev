package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	requestmetrics "github.com/lesi97/lesi.dev/internal/request_metrics"
)

func (db *DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	start := time.Now()
	rows, err := db.Pool.Query(ctx, sql, args...)
	requestmetrics.AddDatabaseCallsDuration(ctx, time.Since(start), err)
	return rows, err
}
