package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	twitch_store "github.com/lesi97/lesi.dev/internal/domains/twitch/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	logger *utils.Logger
	store  twitch_store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB, redis *redis.Client) (*Handler, error) {
	store, err := twitch_store.NewStore(db, logger, redis)
	if err != nil {
		return nil, err
	}

	return &Handler{
		logger: logger,
		store:  store,
	}, nil
}
