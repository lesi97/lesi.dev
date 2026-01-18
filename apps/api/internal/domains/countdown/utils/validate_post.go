package utils

import (
	"errors"

	"github.com/lesi97/lesi.dev/internal/domains/countdown/model"
)

func ValidatePost(data *model.PostData) error {
	
	if data == nil {
		return errors.New("post data is required")
	}

	if data.TargetDate.IsZero() {
		return errors.New("target_date is required")
	}
	
	if data.Message == "" {
		return errors.New("message is required")
	}
	
	if data.FallbackMessage == "" {
		return errors.New("fallback_message is required")
	}
	
	return nil
}
