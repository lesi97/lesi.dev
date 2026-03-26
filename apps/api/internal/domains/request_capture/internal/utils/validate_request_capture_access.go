package utils

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

func ValidateRequestCaptureAccess(r *http.Request) error {
	apiKey := strings.TrimSpace(os.Getenv("REQUEST_CAPTURE_API_KEY"))
	if apiKey == "" {
		return errors.New("request capture api key not configured")
	}

	if r.Header.Get("X-API-Key") != apiKey {
		return errors.New("request capture api key denied")
	}

	return nil
}
