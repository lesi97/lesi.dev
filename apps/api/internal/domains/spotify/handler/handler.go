package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	logger *utils.Logger
	store  store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB, redis *redis.Client) (*Handler, error) {
	musicStore := store.NewStore(db, logger, redis)
	return &Handler{
		logger: logger,
		store:  musicStore,
	}, nil
}
