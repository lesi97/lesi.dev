package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	ani "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/store"
	hu "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/utils/http"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger    *utils.Logger
	store     ani.Methods
	httpUtils hu.Methods
}

func NewHandler(logger *utils.Logger, db *db.DB) (*Handler, error) {

	store, err := ani.NewStore(db, logger)
	if err != nil {
		return nil, err
	}

	httpUtils := hu.NewStore()

	return &Handler{
		logger:    logger,
		store:     store,
		httpUtils: httpUtils,
	}, nil
}
