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
	ApiProcessingMS int64
	FetchCallsMS    int64
	DatabaseCallsMS int64
	NonceElapsedMS  *int64
	StatusCode      *int
}
