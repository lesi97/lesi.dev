package handler

import (
	"github.com/lesi97/lesi.dev/internal/domains/countdown/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger     *utils.Logger
	store	   store.Methods
}

func NewCountdownHandler(logger *utils.Logger, store store.Methods) *Handler {
	return &Handler{
		logger: logger,
		store: store,
	}
}



