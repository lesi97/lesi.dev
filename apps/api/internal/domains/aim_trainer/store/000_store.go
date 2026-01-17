package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/model"
	su "github.com/lesi97/lesi.dev/internal/utils"
)


type Methods interface {
	GetLeaderboard(ctx context.Context) ([]model.LeaderboardRow, error)
	GetUser(ctx context.Context, username string) (*AimTrainerRow, error)
	UpsertUpdateUser(ctx context.Context, in model.UpdateInput) (*model.AimTrainerUpdate, error)
}

type Store struct {
	DB     *db.DB
	Logger *su.Logger
}

func NewStore(db *db.DB, logger *su.Logger) *Store {
	return &Store{
		DB:     db,
		Logger: logger,
	}
}
