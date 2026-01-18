package anilist

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	anilistGET(ctx context.Context, url string) ([]byte, error)
	anilistPOST(ctx context.Context, url string, query string, variables map[string]interface{}) ([]byte, error)
}

type Store struct {
	logger *utils.Logger
	env    *model.AnilistEnv
	db     *db.DB
}

func NewStore(logger *utils.Logger, env *model.AnilistEnv, db *db.DB) *Store {
	return &Store{
		logger: logger,
		env:    env,
		db:     db,
	}
}
