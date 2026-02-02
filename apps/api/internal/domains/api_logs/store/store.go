package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/api_logs/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	InsertApiLog(ctx context.Context, log model.ApiLog) error
}

type Store struct {
	DB     *db.DB
	Logger *utils.Logger
}

func NewStore(db *db.DB, logger *utils.Logger) *Store {
	return &Store{
		DB:     db,
		Logger: logger,
	}
}
