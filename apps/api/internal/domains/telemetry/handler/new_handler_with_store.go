package handler

import (
	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
