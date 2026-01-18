package utils

type bungieSearch struct {
	Response        []bungieSearchResponse `json:"Response"`
	ErrorCode       int                   `json:"ErrorCode"`
	ThrottleSeconds int                   `json:"ThrottleSeconds"`
	ErrorStatus     string                `json:"ErrorStatus"`
	Message         string                `json:"Message"`
	MessageData     struct{}              `json:"MessageData"`
}
