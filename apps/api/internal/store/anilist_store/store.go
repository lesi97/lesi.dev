package anilist_store

import (
	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AnilistStoreInterface interface {}

type AnilistStore struct {
	store.StoreBase
}

func NewStore(db *database.DB, logger *utils.Logger) *AnilistStore {
	return &AnilistStore{
		StoreBase: store.NewStoreBase(db, logger),
	}
}

