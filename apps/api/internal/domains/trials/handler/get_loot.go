package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetLoot(w http.ResponseWriter, r *http.Request) {
	message := h.store.GetLoot()
	utils.TextResponse(w, http.StatusOK, *message)
}
