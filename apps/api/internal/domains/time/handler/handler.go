package handler

import (
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger *utils.Logger
}

func NewHandler(logger *utils.Logger) (*Handler, error) {
	return &Handler{
		logger: logger,
	}, nil
}
