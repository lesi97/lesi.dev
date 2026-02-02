package middleware

import (
	"strconv"
	"time"
)

func GetNonceElapsed(nonce string, now time.Time) (time.Duration, bool) {
	if nonce == "" {
		return 0, false
	}

	nonceMilliseconds, err := strconv.ParseInt(nonce, 10, 64)
	if err != nil {
		return 0, false
	}

	if nonceMilliseconds <= 0 {
		return 0, false
	}

	nonceTime := time.Unix(0, nonceMilliseconds*int64(time.Millisecond))
	elapsed := now.Sub(nonceTime)
	if elapsed < 0 {
		return 0, false
	}

	return elapsed, true
}
