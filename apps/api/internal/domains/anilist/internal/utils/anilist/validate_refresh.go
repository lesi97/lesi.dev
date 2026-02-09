package anilist

import (
	"context"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (s *Store) validateRefresh(ctx context.Context) error {
	expired := utils.IsRefreshTokenExpired(s.env.RefreshExpires)
	if expired {
		api, err := utils.RefreshToken(ctx, "Anilist", s.env.ClientId, s.env.ClientSecret, s.env.RefreshToken, s.env.AuthUrl)
		if err != nil {
			return fmt.Errorf("Failed to refresh Anilist API Details: %v", err)
		}
		s.env.AccessToken = api.AccessToken
		s.env.RefreshToken = api.RefreshToken
		s.env.RefreshExpires = *api.RefreshTokenExpiry
		s.db.UpdateApiDetails(
			ctx,
			"Anilist",
			nil,
			nil,
			&api.AccessToken,
			&api.RefreshToken,
			api.RefreshTokenExpiry,
			nil,
			nil,
			s.logger,
		)
	}

	return nil
}
