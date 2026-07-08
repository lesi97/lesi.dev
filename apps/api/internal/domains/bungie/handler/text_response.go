package handler

import (
	"net/http"
	"strings"

	"github.com/lesi97/lesi.dev/internal/utils"
)

const twitchMessageCharacterLimit = 400

func textResponse(w http.ResponseWriter, status int, message string) {
	message = sanitiseTextResponse(message)
	utils.TruncateString(&message, twitchMessageCharacterLimit)
	utils.TextResponse(w, status, message)
}

func sanitiseTextResponse(message string) string {
	if !strings.Contains(message, "\"") {
		return message
	}

	return strings.ReplaceAll(message, "\"", "")
}
