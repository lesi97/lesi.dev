package middleware

import "net/http"

func GetBotMetadata(headers http.Header) (*string, *string, *string) {
	userDisplayName, channelDisplayName, ok := GetNightbotDisplayNames(headers)
	if ok {
		botType := "nightbot"
		return &botType, &channelDisplayName, &userDisplayName
	}

	streamElementsChannel := headers.Get("X-Streamelements-Channel")
	if streamElementsChannel != "" {
		botType := "streamelements"
		return &botType, &streamElementsChannel, nil
	}

	return nil, nil, nil
}
