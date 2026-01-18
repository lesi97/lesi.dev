package handler

import (
	twitch_store "github.com/lesi97/lesi.dev/internal/domains/twitch/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithStore(logger *utils.Logger, store twitch_store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}
