package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	local_store "github.com/lesi97/lesi.dev/internal/domains/local/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger *utils.Logger
	store  local_store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB) *Handler {
	store := local_store.NewStore(db, logger)

	return &Handler{
		logger: logger,
		store:  store,
	}
}
