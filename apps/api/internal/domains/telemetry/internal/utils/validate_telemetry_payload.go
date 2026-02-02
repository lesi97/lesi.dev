package utils

import (
	"errors"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/model"
)

func ValidateTelemetryPayload(payload model.TelemetryPayload) error {
	if payload.Timestamp == "" {
		return errors.New("timestamp is required")
	}
	if payload.Route == "" {
		return errors.New("route is required")
	}
	if payload.UserAgent == "" {
		return errors.New("userAgent is required")
	}
	if payload.IP != nil && *payload.IP == "" {
		return errors.New("ip must be non-empty when provided")
	}
	if payload.Error != nil && *payload.Error == "" {
		return errors.New("error must be non-empty when provided")
	}

	if _, err := time.Parse(time.RFC3339, payload.Timestamp); err != nil {
		return errors.New("timestamp must be RFC3339")
	}

	return nil
}
