package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/local/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	UpdateApiDetails(ctx context.Context, input model.UpdateApiDetailsInput) error
}

type Store struct {
	DB     *db.DB
	Logger *utils.Logger
}
