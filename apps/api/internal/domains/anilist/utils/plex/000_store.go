package plex

import (
	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Store struct {
	logger           *utils.Logger
	env              *model.PlexEnv
	seasonCountCache map[string]map[int]int
}

func NewStore(logger *utils.Logger, env *model.PlexEnv) *Store {
	return &Store{
		logger:           logger,
		env:              env,
		seasonCountCache: make(map[string]map[int]int),
	}
}
