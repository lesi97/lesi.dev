package model

type PostBody struct {
	Username string   `json:"username"`
	Score    *float64 `json:"score,omitempty"`
	Accuracy *float64 `json:"accuracy,omitempty"`
}

type LeaderboardRow struct {
	Username string  `json:"username"`
	Score    float64 `json:"score"`
	Accuracy float64 `json:"accuracy"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LeaderboardResponse struct {
	Data []LeaderboardRow `json:"data"`
}

type UpdateInput struct {
	Username        string
	Score           *float64
	Accuracy        *float64
	HasNiceTriggers bool
	HasMiloAttacks  bool
	UpdatedAtISO    string
}

type AimTrainerUpdate struct {
	Username        string   `json:"username"`
	UpdatedAt       string   `json:"updated_at"`
	CompletedRounds int64    `json:"completed_rounds"`
	Score           *float64 `json:"score,omitempty"`
	Accuracy        *float64 `json:"accuracy,omitempty"`
	NiceTriggers    *int64   `json:"nice_triggers,omitempty"`
	MiloAttacks     *int64   `json:"milo_attacks,omitempty"`
}

type ExistingRow struct {
	Score           float64
	Accuracy        float64
	CompletedRounds int64
	MiloAttacks     int64
	NiceTriggers    int64
}
