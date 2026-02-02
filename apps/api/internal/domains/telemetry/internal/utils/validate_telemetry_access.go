package utils

import (
	"errors"
	"net/http"
	"os"
)

func ValidateTelemetryAccess(r *http.Request) error {
	allowedOrigin := os.Getenv("WEB_URL")
	if allowedOrigin == "" {
		return errors.New("telemetry origin not configured")
	}

	origin, ok := GetRequestOrigin(r)
	if !ok || origin != allowedOrigin {
		return errors.New("telemetry origin denied")
	}

	apiKey := os.Getenv("TELEMETRY_API_KEY")
	if apiKey == "" {
		return errors.New("telemetry api key not configured")
	}

	if r.Header.Get("X-Telemetry-Key") != apiKey {
		return errors.New("telemetry api key denied")
	}

	return nil
}
