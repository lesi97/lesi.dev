package store

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type spotifyCredentials struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AccountsURL  string
	APIURL       string
}

func (s *Store) spotifyCredentials(ctx context.Context) (spotifyCredentials, error) {
	creds := spotifyCredentials{
		ClientID:     stringValue(s.spotify.client_id),
		ClientSecret: stringValue(s.spotify.client_secret),
		RedirectURL:  strings.TrimSpace(os.Getenv("SPOTIFY_REDIRECT_URL")),
		AccountsURL:  s.spotify.base_url,
		APIURL:       "https://api.spotify.com/v1",
	}

	if creds.RedirectURL == "" && s.apiUrl != nil {
		creds.RedirectURL = fmt.Sprintf("%v/v1/auth/spotify/callback", *s.apiUrl)
	}

	if s.DB != nil {
		apiDetails, err := s.DB.FetchApiDetails(ctx, "Spotify", s.Logger)
		if err != nil {
			return spotifyCredentials{}, err
		}
		if apiDetails != nil {
			if apiDetails.ClientID != nil && strings.TrimSpace(*apiDetails.ClientID) != "" {
				creds.ClientID = strings.TrimSpace(*apiDetails.ClientID)
			}
			if apiDetails.ClientSecret != nil && strings.TrimSpace(*apiDetails.ClientSecret) != "" {
				creds.ClientSecret = strings.TrimSpace(*apiDetails.ClientSecret)
			}
			if apiDetails.RedirectURL != nil && strings.TrimSpace(*apiDetails.RedirectURL) != "" {
				creds.RedirectURL = strings.TrimSpace(*apiDetails.RedirectURL)
			}
			if apiDetails.BaseURL != nil && strings.TrimSpace(*apiDetails.BaseURL) != "" {
				creds.APIURL = strings.TrimSpace(*apiDetails.BaseURL)
			}
		}
	}

	if creds.ClientID == "" {
		return spotifyCredentials{}, fmt.Errorf("missing Spotify client id")
	}
	if creds.ClientSecret == "" {
		return spotifyCredentials{}, fmt.Errorf("missing Spotify client secret")
	}
	if creds.RedirectURL == "" {
		return spotifyCredentials{}, fmt.Errorf("missing Spotify redirect url")
	}

	return creds, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
