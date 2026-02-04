package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetPlayerCount(w http.ResponseWriter, r *http.Request) {
	
	message, err := h.store.GetPlayerCount(r.Context())
	if err != nil {
		h.logger.Printf("ERROR: GetPlayerCount | %v", err.Error())
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}
