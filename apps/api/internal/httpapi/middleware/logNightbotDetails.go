package middleware

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func LogNightbotDetails(logger *utils.Logger, header http.Header, responseBody string, status int) {
	userDisplayName, channelDisplayName, ok := GetNightbotDisplayNames(header)
	if !ok || status == http.StatusNotFound {
		return
	}

	reset := utils.Colours["reset"]
	nightbotColour := utils.Colours["brightMagenta"]
	logger.Printf("%vNightbot User: %v%v", nightbotColour, userDisplayName, reset)
	logger.Printf("%vNightbot Channel: %v%v", nightbotColour, channelDisplayName, reset)
	logger.Printf("%vNightbot Response: %v%v", nightbotColour, responseBody, reset)
}
