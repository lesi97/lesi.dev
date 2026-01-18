package model

type TwitchChatters struct {
	Data       []TwitchUser     `json:"data"`
	Pagination TwitchPagination `json:"pagination"`
	Total      int             `json:"total"`
}
