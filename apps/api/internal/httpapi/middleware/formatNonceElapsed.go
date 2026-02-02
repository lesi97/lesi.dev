package middleware

import (
	"strings"
	"time"
)

func FormatNonceElapsed(elapsed time.Duration) string {
	raw := elapsed.String()
	if elapsed < time.Minute {
		return raw
	}

	minuteIndex := strings.Index(raw, "m")
	if minuteIndex == -1 || minuteIndex == len(raw)-1 {
		return raw
	}

	return raw[:minuteIndex+1] + " " + raw[minuteIndex+1:]
}
