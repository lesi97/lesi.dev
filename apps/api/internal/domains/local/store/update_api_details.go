package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/domains/local/model"
)

func (s *Store) UpdateApiDetails(ctx context.Context, input model.UpdateApiDetailsInput) error {
	return s.DB.UpdateApiDetails(
		ctx,
		input.Name,
		&input.ClientID,
		&input.ClientSecret,
		&input.AccessToken,
		&input.RefreshToken,
		&input.RefreshTokenExpiry,
		&input.BaseURL,
		&input.RedirectURL,
		s.Logger,
	)
}
