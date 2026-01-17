package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) getRandomTarot(w http.ResponseWriter, r *http.Request) {
	tarot := h.store.GetRandomTarot()
	utils.TextResponse(w, http.StatusOK, *tarot)
}