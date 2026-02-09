package store

import (
	"context"
	"fmt"

	oauth_utils "github.com/lesi97/lesi.dev/internal/domains/auth/internal/utils/oauth"
)

func (s *Store) TwitchModCallback(ctx context.Context, code string) error {
	cfg := oauth_utils.ProviderConfig{
		Application:  "Twitch_GO",
		TokenURL:     fmt.Sprintf("%v/token", s.twitch.base_url),
		ClientID:     s.twitch.client_id,
		ClientSecret: s.twitch.client_secret,
		RedirectURL:  fmt.Sprintf("%v/v1/auth/twitch/callback", *s.apiUrl),
	}

	return oauth_utils.OAuthCallback(ctx, s.DB, s.Logger, code, cfg)
}
