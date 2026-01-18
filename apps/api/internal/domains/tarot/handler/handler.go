package handler

import (
	"github.com/lesi97/lesi.dev/internal/domains/tarot/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger *utils.Logger
	store  store.TarotStoreInterface
}

func NewHandler(logger *utils.Logger) (*Handler, error) {
	store := store.NewStore(logger)
	return &Handler{
		logger: logger,
		store:  store,
	}, nil
}
