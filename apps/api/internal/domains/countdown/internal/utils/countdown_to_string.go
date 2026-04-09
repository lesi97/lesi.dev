package utils

import (
	"strings"
	"time"
)

/*
Utility function to convert check if the date from the countdown api route has been reached or not

Returns string "Passed" if the date has passed

Returns string difference in time if the countdown has not yet been reached
*/
func CountdownToString(target time.Time) string {
	now := time.Now()
	if now.After(target) {
		return "Passed"
	}

	diff := target.Sub(now).Round(time.Second)

	days := int(diff.Hours()) / 24
	hours := int(diff.Hours()) % 24
	minutes := int(diff.Minutes()) % 60
	seconds := int(diff.Seconds()) % 60

	timeParts := []string{}

	if days > 0 {
		timeParts = append(timeParts, CountdownTimePartToString(days, "day"))
	}
	if days > 0 || hours > 0 {
		timeParts = append(timeParts, CountdownTimePartToString(hours, "hour"))
	}
	if days > 0 || hours > 0 || minutes > 0 {
		timeParts = append(timeParts, CountdownTimePartToString(minutes, "minute"))
	}
	timeParts = append(timeParts, CountdownTimePartToString(seconds, "second"))

	return strings.Join(timeParts, ", ")
}
