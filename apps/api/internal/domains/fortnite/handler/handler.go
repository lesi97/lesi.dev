package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/db"
	fortnite_store "github.com/lesi97/lesi.dev/internal/domains/fortnite/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	logger *utils.Logger
	store  fortnite_store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB, redis *redis.Client, httpClient *http.Client) (*Handler, error) {
	store, err := fortnite_store.NewStore(db, logger, redis, httpClient)
	if err != nil {
		return nil, err
	}

	return &Handler{
		logger: logger,
		store:  store,
	}, nil
}
