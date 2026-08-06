package handler

import (
	"net/http"

	shared_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) GetLatestScrobble(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

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
