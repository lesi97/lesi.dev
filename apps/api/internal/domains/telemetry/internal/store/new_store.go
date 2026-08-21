package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	InsertTelemetry(ctx context.Context, payload model.TelemetryPayload) error
	UpsertDateMemeClick(ctx context.Context, input model.DateMemeClickInput) error
}

type Store struct {
	DB     *db.DB
	Logger *utils.Logger
}

func NewStore(db *db.DB, logger *utils.Logger) *Store {
	return &Store{
		DB:     db,
		Logger: logger,
	}
}
