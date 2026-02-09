package store

import (
	"context"
	"fmt"

	oauth_utils "github.com/lesi97/lesi.dev/internal/domains/auth/internal/utils/oauth"
)

func (s *Store) AnilistCallback(ctx context.Context, code string) error {
	cfg := oauth_utils.ProviderConfig{
		Application:  "Anilist",
		TokenURL:     fmt.Sprintf("%v/token", s.anilist.base_url),
		ClientID:     s.anilist.client_id,
		ClientSecret: s.anilist.client_secret,
		RedirectURL:  fmt.Sprintf("%v/v1/auth/anilist/callback", *s.apiUrl),
	}

	return oauth_utils.OAuthCallback(ctx, s.DB, s.Logger, code, cfg)
}
