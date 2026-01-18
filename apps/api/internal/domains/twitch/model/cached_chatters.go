package model

import "time"

type CachedChatters struct {
	FetchedAt time.Time      `json:"fetched_at"`
	Data      TwitchChatters `json:"data"`
}
