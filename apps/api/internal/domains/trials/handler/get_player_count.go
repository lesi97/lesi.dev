package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetPlayerCount(w http.ResponseWriter, r *http.Request) {
	message := h.store.GetPlayerCount(r.Context())
	utils.TextResponse(w, http.StatusOK, *message)
}
