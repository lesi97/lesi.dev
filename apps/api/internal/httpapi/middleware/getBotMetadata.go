package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

var fetchStreamElementsChannelDisplayName = FetchStreamElementsChannelDisplayName

func GetBotMetadata(headers http.Header, query url.Values) (*string, *string, *string) {
	userDisplayName, channelDisplayName, ok := GetNightbotDisplayNames(headers)
	if ok {
		botType := "nightbot"
		return &botType, &channelDisplayName, &userDisplayName
	}

	streamElementsChannel := headers.Get("X-Streamelements-Channel")
	if streamElementsChannel != "" {
		botType := "streamelements"
		userName := strings.TrimSpace(query.Get("sender"))
		var user *string
		if userName != "" {
			user = &userName
		}

		channelDisplayName, ok := fetchStreamElementsChannelDisplayName(streamElementsChannel)
		if ok {
			return &botType, &channelDisplayName, user
		}

		return &botType, &streamElementsChannel, user
	}

	return nil, nil, nil
}
