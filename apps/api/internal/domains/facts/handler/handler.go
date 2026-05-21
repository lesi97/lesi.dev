package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	facts_store "github.com/lesi97/lesi.dev/internal/domains/facts/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger *utils.Logger
	store  facts_store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB) (*Handler, error) {
	store, err := facts_store.NewStore(db, logger)
	if err != nil {
		return nil, err
	}

	return &Handler{
		logger: logger,
		store:  store,
	}, nil
}
