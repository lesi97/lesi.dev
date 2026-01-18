package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/countdown/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	GetCountdownByID(ctx context.Context, id string) (*string, error)
	InsertCountdown(ctx context.Context, data model.PostData) (*string, error)
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
