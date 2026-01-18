package store

import (
	"context"
	"fmt"

	twitch_utils "github.com/lesi97/lesi.dev/internal/domains/auth/utils/twitch"
)

func (s *Store) TwitchModCallback(code string) error {
	cfg := twitch_utils.OAuthProviderConfig{
		Application:  "Twitch_GO",
		TokenURL:     fmt.Sprintf("%v/token", s.twitch.base_url),
		ClientID:     s.twitch.client_id,
		ClientSecret: s.twitch.client_secret,
		RedirectURL:  fmt.Sprintf("%v/v1/auth/twitch/callback", *s.apiUrl),
	}

	return twitch_utils.OAuthCallback(context.Background(), s.DB, s.Logger, code, cfg)
}
