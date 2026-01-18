package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/domains/countdown/model"
)

func (s *Store) InsertCountdown(ctx context.Context, data model.PostData) (*string, error) {
	return data.Insert(s.DB, ctx)
}
