package model

type TrialsPlatformRecentStats struct {
	HourEnding  string `json:"hourEnding"`
	PlayerCount int64  `json:"playerCount"`
	MatchCount  int    `json:"matchCount"`
	UpdatedAt   string `json:"updatedAt"`
}
