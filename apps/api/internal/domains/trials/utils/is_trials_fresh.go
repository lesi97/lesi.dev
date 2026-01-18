package utils

import "time"

func IsTrialsFresh(updatedAt time.Time) bool {
	return time.Since(updatedAt) < 90*time.Minute
}
