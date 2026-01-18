package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	trials_store "github.com/lesi97/lesi.dev/internal/domains/trials/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger *utils.Logger
	store  trials_store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB) *Handler {
	store := trials_store.NewStore(db, logger)

	return &Handler{
		logger: logger,
		store:  store,
	}
}
