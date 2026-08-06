package utils

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

func ValidateScrobbleAccess(r *http.Request) error {
	apiKey := strings.TrimSpace(os.Getenv("SCROBBLE_API_KEY"))
	if apiKey == "" {
		return errors.New("scrobble api key not configured")
	}

	if r.Header.Get("X-API-Key") != apiKey {
		return errors.New("scrobble api key denied")
	}

	return nil
}
