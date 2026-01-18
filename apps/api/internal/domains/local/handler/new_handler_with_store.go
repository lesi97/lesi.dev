package handler

import (
	local_store "github.com/lesi97/lesi.dev/internal/domains/local/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store local_store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
