package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/api/internal/store"
	"github.com/lesi97/lesi.dev/api/internal/utils"
)

type CountdownHandler struct {
	logger     *utils.Logger
	countdownStore store.CountdownStore
}

func NewCountdownHandler(logger *utils.Logger, countdownStore store.CountdownStore)  *CountdownHandler {
	return &CountdownHandler{
		logger: logger,
		countdownStore: countdownStore,
	}
}

func (h *CountdownHandler) HandleGetCountdown(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("ERROR: handleGetCountdown: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}

	countdownMessage, err := h.countdownStore.GetCountdownByID(r.Context(), idParam)
	if err != nil {
		h.logger.Printf("ERROR: GetCountdownByID: %v", err)
		utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.TextResponse(w, http.StatusOK, *countdownMessage)
}