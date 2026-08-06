package store

import (
	"context"
	"fmt"
	"net/url"
)

const spotifyRecentlyPlayedScope = "user-read-recently-played"

func (s *Store) SpotifyAuthUrl(ctx context.Context) (*string, error) {
	creds, err := s.spotifyCredentials(ctx)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize", creds.AccountsURL))
	if err != nil {
		return nil, err
	}

	query := baseURL.Query()
	query.Set("client_id", creds.ClientID)
	query.Set("redirect_uri", creds.RedirectURL)
	query.Set("response_type", "code")
	query.Set("scope", spotifyRecentlyPlayedScope)

	baseURL.RawQuery = query.Encode()
	str := baseURL.String()
	return &str, nil
}
