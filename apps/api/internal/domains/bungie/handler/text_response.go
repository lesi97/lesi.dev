package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

const twitchMessageCharacterLimit = 400

func textResponse(w http.ResponseWriter, status int, message string) {
	utils.TruncateString(&message, twitchMessageCharacterLimit)
	utils.TextResponse(w, status, message)
}
