package utils

import "time"

func IsRefreshTokenExpired(expiryTime int64) bool {
	return time.Now().UnixMilli() > expiryTime
}
