package handler

import (
	steam_store "github.com/lesi97/lesi.dev/internal/domains/steam/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store steam_store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
