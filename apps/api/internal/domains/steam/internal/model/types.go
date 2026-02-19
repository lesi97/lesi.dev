package model

type SteamCurrentPlayersData struct {
	Response SteamCurrentPlayersResponse `json:"response"`
}

type SteamCurrentPlayersResponse struct {
	PlayerCount int64 `json:"player_count"`
	Result      int64 `json:"result"`
}

type SteamApp struct {
	AppID int64  `json:"appid"`
	Name  string `json:"name"`
}

type SteamAppList struct {
	Apps []SteamApp `json:"apps"`
}

type SteamAppListData struct {
	AppList SteamAppList `json:"applist"`
}

type SteamStoreSearchData struct {
	Items []SteamStoreSearchItem `json:"items"`
}

type SteamStoreSearchItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SteamAppDetailsByIDData map[string]SteamAppDetailsData

type SteamAppDetailsData struct {
	Success bool                `json:"success"`
	Data    SteamAppDetailsItem `json:"data"`
}

type SteamAppDetailsItem struct {
	Name string `json:"name"`
}
