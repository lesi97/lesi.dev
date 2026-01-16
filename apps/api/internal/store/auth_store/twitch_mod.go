package auth_store

import (
	"context"
	"fmt"
	"net/url"
)

func (s *AuthStore) TwitchModCallback(code string) error {

	cfg := oauthProviderConfig{
		Application: "Twitch_GO",
		TokenURL: fmt.Sprintf("%v/token", s.twitch.base_url),
		ClientID: s.twitch.client_id,
		ClientSecret: s.twitch.client_secret,
		RedirectURL: fmt.Sprintf("%v/v1/auth/twitch/callback", *s.apiUrl),
	}

	return s.oauthCallback(context.Background(), code, cfg)
}


func (s *AuthStore) TwitchModAuthUrl() (*string, error) {

	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize", s.twitch.base_url))
	if err != nil {
		return nil, err
	}
	query := baseURL.Query()
	query.Set("client_id", *s.twitch.client_id)
	query.Set("redirect_uri", fmt.Sprintf("%v/v1/auth/twitch/callback", *s.apiUrl))
	query.Set("response_type", "code")
	query.Set("scope", "moderator:read:chatters")

	baseURL.RawQuery = query.Encode()
	str := baseURL.String()
	return &str, nil
}