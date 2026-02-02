package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger *utils.Logger
	store  store.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB) (*Handler, error) {
	telemetryStore := store.NewStore(db, logger)
	return &Handler{
		logger: logger,
		store:  telemetryStore,
	}, nil
}
