package store

import (
	"context"
	"fmt"

	anilist_utils "github.com/lesi97/lesi.dev/internal/domains/auth/utils/anilist"
)

func (s *Store) AnilistCallback(code string) error {
	cfg := anilist_utils.OAuthProviderConfig{
		Application:  "Anilist",
		TokenURL:     fmt.Sprintf("%v/token", s.anilist.base_url),
		ClientID:     s.anilist.client_id,
		ClientSecret: s.anilist.client_secret,
		RedirectURL:  fmt.Sprintf("%v/v1/auth/anilist/callback", *s.apiUrl),
	}

	return anilist_utils.OAuthCallback(context.Background(), s.DB, s.Logger, code, cfg)
}
