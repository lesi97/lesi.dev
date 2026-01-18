package model

type SteamResponse struct {
	PlayerCount int64 `json:"player_count"`
	Result      int   `json:"result"`
}
