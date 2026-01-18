package handler

import (
	"github.com/lesi97/lesi.dev/internal/store/auth_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger         *utils.Logger
	store 			auth_store.AuthStoreInterface
}

func NewHandler(logger *utils.Logger, store auth_store.AuthStoreInterface) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}

