package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	auth_store "github.com/lesi97/lesi.dev/internal/domains/auth/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger         *utils.Logger
	store          auth_store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB) (*Handler, error) {
	store, err := auth_store.NewStore(db, logger)
	if err != nil {
		return nil, err
	}

	return &Handler{
		logger: logger,
		store:  store,
	}, nil
}

