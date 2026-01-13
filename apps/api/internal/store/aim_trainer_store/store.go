package aim_trainer_store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)


type AimTrainerStoreInterface interface {
	GetLeaderboard(ctx context.Context) ([]LeaderboardRow, error)
	GetUser(ctx context.Context, username string) (*AimTrainerRow, error)
	UpsertUpdateUser(ctx context.Context, in UpdateInput) (*AimTrainerUpdate, error)
}

type AimTrainerStore struct {
	store.StoreBase
}

func NewStore(db *database.DB, logger *utils.Logger) *AimTrainerStore {
	return &AimTrainerStore{
		StoreBase: store.NewStoreBase(db, logger),
	}
}
