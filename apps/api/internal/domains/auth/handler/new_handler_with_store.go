package handler

import (
	auth_store "github.com/lesi97/lesi.dev/internal/domains/auth/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store auth_store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
