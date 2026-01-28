package middleware

import (
	"net/http"
	"net/url"
)

func GetNightbotDisplayNames(headers http.Header) (string, string, bool) {
	userHeader := headers.Get("Nightbot-User")
	if userHeader == "" {
		return "", "", false
	}

	channelHeader := headers.Get("Nightbot-Channel")
	if channelHeader == "" {
		return "", "", false
	}

	userValues, err := url.ParseQuery(userHeader)
	if err != nil {
		return "", "", false
	}

	channelValues, err := url.ParseQuery(channelHeader)
	if err != nil {
		return "", "", false
	}

	userDisplayName := userValues.Get("displayName")
	if userDisplayName == "" {
		return "", "", false
	}

	channelDisplayName := channelValues.Get("displayName")
	if channelDisplayName == "" {
		return "", "", false
	}

	return userDisplayName, channelDisplayName, true
}
