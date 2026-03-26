package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/request_capture/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	InsertRequestCapture(ctx context.Context, input model.RequestCaptureInput) (int64, error)
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
