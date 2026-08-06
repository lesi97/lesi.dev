package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Store) GetLatestPlayedText(ctx context.Context) (*string, error) {
	var text string
	err := s.DB.QueryRow(ctx, `
		SELECT a.name || ' - ' || t.name
		FROM music.scrobbles s
		JOIN music.artists a ON a.id = s.artist_id
		JOIN music.tracks t ON t.id = s.track_id
		ORDER BY s.scrobbled_at DESC, s.id DESC
		LIMIT 1
	`).Scan(&text)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &text, nil
}
