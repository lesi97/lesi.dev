package store

import (
	"fmt"
	"net/url"
)

func (s *Store) TwitchModAuthUrl() (*string, error) {
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
