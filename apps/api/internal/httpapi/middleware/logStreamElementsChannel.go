package middleware

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func LogStreamElementsChannel(logger *utils.Logger, headers http.Header, responseBody string) {
	channelId := headers.Get("X-Streamelements-Channel")
	if channelId == "" {
		return
	}

	go func() {
		displayName, ok := FetchStreamElementsChannelDisplayName(channelId)
		if !ok {
			return
		}
		logger.Printf(
			"%vStreamElements Channel: %v%v",
			utils.Colours["brightMagenta"],
			displayName,
			utils.Colours["reset"],
		)
		logger.Printf(
			"%vResponse: %v%v",
			utils.Colours["brightMagenta"],
			responseBody,
			utils.Colours["reset"],
		)
	}()
}
