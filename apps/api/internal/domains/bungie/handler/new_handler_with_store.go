package handler

import (
	bungie_store "github.com/lesi97/lesi.dev/internal/domains/bungie/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store bungie_store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
