package utils

import (
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
)

func ParseTrialsUpdatedAt(data *model.TrialsData) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, data.Platforms.Num0.RecentStats.UpdatedAt)
}
