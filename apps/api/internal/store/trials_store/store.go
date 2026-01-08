package trials_store

import (
	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type TrialsStoreInterface interface {
	GetLoot() *string
	GetPlayerCount() *string
}

type TrialsStore struct {
	store.StoreBase
}

func NewStore(db *database.DB, logger *utils.Logger) *TrialsStore {
	return &TrialsStore{
		StoreBase: store.NewStoreBase(db, logger),
	}
}


