package store

import (
	"net/url"
	"strings"
)

func DecodeGameName(gameName string) string {
	trimmedGameName := strings.TrimSpace(gameName)
	if trimmedGameName == "" {
		return ""
	}

	decodedGameName, err := url.QueryUnescape(trimmedGameName)
	if err != nil {
		return trimmedGameName
	}

	if strings.TrimSpace(decodedGameName) == "" {
		return trimmedGameName
	}

	return decodedGameName
}
