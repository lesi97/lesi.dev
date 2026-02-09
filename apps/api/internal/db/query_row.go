package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	requestmetrics "github.com/lesi97/lesi.dev/internal/request_metrics"
)

func (db *DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	start := time.Now()
	row := db.Pool.QueryRow(ctx, sql, args...)
	requestmetrics.AddDatabaseCallsDuration(ctx, time.Since(start), nil)
	return row
}
