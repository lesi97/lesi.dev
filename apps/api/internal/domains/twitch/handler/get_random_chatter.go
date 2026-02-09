package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetRandomChatter(w http.ResponseWriter, r *http.Request) {
	streamerName := chi.URLParam(r, "streamer")
	if streamerName == "" {
		utils.TextResponse(w, http.StatusOK, "You must declare a streamer for this to work")
		return
	}

	username, err := h.store.RandomViewer(r.Context(), streamerName)
	if err != nil {
		h.logger.Errorf("ERROR: %v", err)
		utils.TextResponse(w, http.StatusOK, "An error has occurred, take note of the current time and tell Lesi when this happened")
		return
	}
	utils.TextResponse(w, http.StatusOK, *username)
}
