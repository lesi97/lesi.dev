package utils

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

func ValidateTelemetryAccess(r *http.Request) error {
	origin, ok := GetRequestOrigin(r)
	if !ok {
		return errors.New("telemetry origin denied")
	}

	webURL := strings.TrimSpace(os.Getenv("WEB_URL"))
	vercelURL := strings.TrimSpace(os.Getenv("VERCEL_URL"))

	if webURL == "" {
		return errors.New("telemetry origin not configured")
	}

	isAllowed := false

	if origin == webURL {
		isAllowed = true
	}

	if vercelURL != "" && origin == vercelURL {
		isAllowed = true
	}

	if !isAllowed {
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