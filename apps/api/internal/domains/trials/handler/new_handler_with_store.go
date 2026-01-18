package handler

import (
	trials_store "github.com/lesi97/lesi.dev/internal/domains/trials/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store trials_store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
