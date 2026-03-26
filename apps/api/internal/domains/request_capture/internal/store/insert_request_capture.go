package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/domains/request_capture/internal/model"
)

func (s *Store) InsertRequestCapture(ctx context.Context, input model.RequestCaptureInput) (int64, error) {
	query := `
		INSERT INTO public.request_capture (
			path,
			ip,
			user_agent,
			content_type,
			content_length,
			headers,
			body
		)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7::jsonb)
		RETURNING id
	`

	var id int64
	err := s.DB.QueryRow(
		ctx,
		query,
		input.Path,
		input.IP,
		input.UserAgent,
		input.ContentType,
		input.ContentLength,
		input.Headers,
		input.Body,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
