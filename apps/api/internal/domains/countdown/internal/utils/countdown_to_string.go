package utils

import (
	"fmt"
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

	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)
}
