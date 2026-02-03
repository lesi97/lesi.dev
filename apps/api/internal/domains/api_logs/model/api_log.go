package model

import "time"

type ApiLog struct {
	Timestamp       time.Time
	Route           string
	IP              string
	Channel         *string
	User            *string
	BotType         *string
	Response        string
	ExecutionTimeMS int64
	NonceElapsedMS  *int64
	StatusCode		*int
}
