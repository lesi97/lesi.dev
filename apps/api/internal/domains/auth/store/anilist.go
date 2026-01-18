package store

import (
	"context"
	"fmt"
	"net/url"
)

func (s *AuthStore) AnilistCallback(code string) error {

	cfg := oauthProviderConfig{
		Application: "Anilist",
		TokenURL: fmt.Sprintf("%v/token", s.anilist.base_url),
		ClientID: s.anilist.client_id,
		ClientSecret: s.anilist.client_secret,
		RedirectURL: fmt.Sprintf("%v/v1/auth/anilist/callback", *s.apiUrl),
	}

	return s.oauthCallback(context.Background(), code, cfg)
}


func (s *AuthStore) AnilistAuthUrl() (*string, error) {

	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize", s.anilist.base_url))
	if err != nil {
		return nil, err
	}

	query := baseURL.Query()
	query.Set("client_id", *s.anilist.client_id)
	query.Set("redirect_uri", fmt.Sprintf("%v/v1/auth/anilist/callback", *s.apiUrl))
	query.Set("response_type", "code")

	baseURL.RawQuery = query.Encode()
	str := baseURL.String()
	return &str, nil
}