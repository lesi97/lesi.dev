package store

import (
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewStore(db *db.DB, logger *utils.Logger) *Store {
	return &Store{
		DB:     db,
		Logger: logger,
	}
}
