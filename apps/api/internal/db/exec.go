package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	requestmetrics "github.com/lesi97/lesi.dev/internal/request_metrics"
)

func (db *DB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	start := time.Now()
	commandTag, err := db.Pool.Exec(ctx, sql, arguments...)
	requestmetrics.AddDatabaseCallsDuration(ctx, time.Since(start), err)
	return commandTag, err
}
