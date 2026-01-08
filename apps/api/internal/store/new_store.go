package store

import (
	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type StoreBase struct {
	DB     *database.DB
	Logger *utils.Logger
}

func NewStoreBase(db *database.DB, logger *utils.Logger) StoreBase {
	return StoreBase{
		DB:     db,
		Logger: logger,
	}
}
