package utils

import "time"

func IsTrialsReportAvailable(now time.Time) bool {
	utc := now.UTC()
	switch utc.Weekday() {
	case time.Friday:
		return utc.Hour() >= 17
	case time.Saturday, time.Sunday, time.Monday:
		return true
	case time.Tuesday:
		return utc.Hour() < 17
	default:
		return false
	}
}
