package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetPlayerCount(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("gameId")
	gameName := r.URL.Query().Get("gameName")

	message, err := h.store.GetPlayerCount(r.Context(), gameID, gameName)
	if err != nil {
		h.logger.Printf("ERROR: GetPlayerCount | %v", err.Error())
		// have to return ok for nightbot to return my error response
		utils.TextResponse(w, http.StatusOK, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}
