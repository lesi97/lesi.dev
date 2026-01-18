package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	bungie_store "github.com/lesi97/lesi.dev/internal/domains/bungie/store"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	logger *utils.Logger
	store  bungie_store.Methods
}

const bungieContextKey = "bungie"

func NewHandler(logger *utils.Logger, db *db.DB, redis *redis.Client) (*Handler, error) {
	store, err := bungie_store.NewStore(db, logger, redis)
	if err != nil {
		return nil, err
	}

	return &Handler{
		logger: logger,
		store:  store,
	}, nil
}
