package handler

import (
	"github.com/lesi97/lesi.dev/internal/domains/countdown/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type handler struct {
	logger     *utils.Logger
	store	   store.Methods
}

func NewCountdownHandler(logger *utils.Logger, store store.Methods)  *handler {
	
	return &handler{
		logger: logger,
		store: store,
	}
}



