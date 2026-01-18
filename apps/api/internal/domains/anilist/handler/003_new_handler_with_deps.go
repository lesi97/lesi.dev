package handler

import (
	ani "github.com/lesi97/lesi.dev/internal/domains/anilist/store"
	hu "github.com/lesi97/lesi.dev/internal/domains/anilist/utils/http"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewHandlerWithDeps(logger *utils.Logger, store ani.Methods, httpUtils hu.Methods) *Handler {
	return &Handler{
		logger:    logger,
		store:     store,
		httpUtils: httpUtils,
	}
}
