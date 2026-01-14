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
		ClientID: s.twitch.api_details.ClientID,
		ClientSecret: s.twitch.api_details.ClientSecret,
		RedirectURL: s.twitch.api_details.RedirectURL,
	}

	return s.oauthCallback(context.Background(), code, cfg)
}


func (s *AuthStore) TwitchModAuthUrl() (*string, error) {

	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize", s.twitch.base_url))
	if err != nil {
		return nil, err
	}
	query := baseURL.Query()
	query.Set("client_id", *s.twitch.api_details.ClientID)
	query.Set("redirect_uri", *s.twitch.api_details.RedirectURL)
	query.Set("response_type", "code")
	query.Set("scope", "moderator:read:chatters")

	baseURL.RawQuery = query.Encode()
	str := baseURL.String()
	return &str, nil
}