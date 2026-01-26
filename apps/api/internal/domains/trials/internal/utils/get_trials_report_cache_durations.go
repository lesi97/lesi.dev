package utils

import (
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
)

func GetTrialsReportCacheDurations(now time.Time, data *model.TrialsData) (time.Duration, time.Duration) {
	if data == nil {
		return 0, 0
	}

	updatedAt, err := ParseTrialsUpdatedAt(data)
	if err != nil {
		return 0, 0
	}

	elapsed := now.Sub(updatedAt)
	if elapsed < 0 {
		elapsed = 0
	}

	remaining := (30 * time.Minute) - elapsed
	if remaining <= 0 {
		return 0, 0
	}

	staleExtension := 5 * time.Minute
	return remaining, remaining + staleExtension
}
