package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) getCountdown(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}

	countdownMessage, err := h.store.GetCountdownByID(r.Context(), idParam)
	if err != nil {
		h.logger.Printf("ERROR: GetCountdownByID: %v", err)
		utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.TextResponse(w, http.StatusOK, *countdownMessage)
}
