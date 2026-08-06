package handler

import (
	"net/http"

	shared_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) GetLatestScrobble(w http.ResponseWriter, r *http.Request) {
	text, err := h.store.GetLatestPlayedText(r.Context())
	if err != nil {
		h.logger.Printf("ERROR: Latest scrobble GET: %v", err)
		shared_utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if text == nil {
		shared_utils.TextResponse(w, http.StatusNotFound, "no scrobbles found")
		return
	}

	shared_utils.TextResponse(w, http.StatusOK, *text)
}
