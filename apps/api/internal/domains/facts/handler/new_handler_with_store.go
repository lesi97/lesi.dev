package handler

import (
	facts_store "github.com/lesi97/lesi.dev/internal/domains/facts/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store facts_store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
