package utils

import (
	"errors"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/model"
)

func ValidateDateMemeClickPayload(payload model.DateMemeClickPayload) error {
	if payload.Route == "" {
		return errors.New("route is required")
	}
	if len(payload.Route) > 256 {
		return errors.New("route must be 256 characters or fewer")
	}
	isDateSlugRoute := strings.HasPrefix(payload.Route, "/date/") && strings.TrimPrefix(payload.Route, "/date/") != ""
	if payload.Route != "/audrey" && !isDateSlugRoute {
		return errors.New("route must be /audrey or /date/{slug}")
	}
	if payload.Action != model.DateMemeClickActionYes && payload.Action != model.DateMemeClickActionNo {
		return errors.New("action must be yes or no")
	}
	if payload.SecretEnding && payload.Action != model.DateMemeClickActionYes {
		return errors.New("secret ending can only be tracked for yes clicks")
	}

	return nil
}
