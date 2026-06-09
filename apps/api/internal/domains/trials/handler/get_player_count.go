package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetPlayerCount(w http.ResponseWriter, r *http.Request) {
	message := h.store.GetPlayerCount(r.Context())
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	utils.TextResponse(w, http.StatusOK, *message)
}
