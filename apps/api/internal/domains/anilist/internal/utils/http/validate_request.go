package utils

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

func (s *Store) ValidateRequest(r *http.Request) (int, string, error) {
	secret := r.URL.Query().Get("secret")

	plexUsername := os.Getenv("PLEX_USERNAME")
	if plexUsername == "" {
		return http.StatusInternalServerError, "", errors.New("PLEX_USERNAME not found in env")
	}

	if secret == "" || secret != os.Getenv("PLEX_WEBHOOK_SECRET") {
		return http.StatusForbidden, plexUsername, errors.New("forbidden")
	}

	ua := r.Header.Get("User-Agent")
	if !strings.Contains(ua, "PlexMediaServer") {
		return http.StatusForbidden, plexUsername, errors.New("forbidden")
	}

	return 0, plexUsername, nil
}
