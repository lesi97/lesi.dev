package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	GetRandomFact(ctx context.Context) (*string, error)
	GetRandomCheeseFact(ctx context.Context) (*string, error)
}

type Store struct {
	DB     *db.DB
	Logger *utils.Logger
}

func NewStore(db *db.DB, logger *utils.Logger) (*Store, error) {
	store := &Store{
		DB:     db,
		Logger: logger,
	}

	return store, nil
}
