package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/api.lesi.dev/internal/store"
	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type CountdownHandler struct {
	logger     *log.Logger
	countdownStore store.CountdownStore
}

func NewCountdownHandler(logger *log.Logger, countdownStore store.CountdownStore)  *CountdownHandler {
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

	countdownData, err := h.countdownStore.GetCountdownByID(idParam)
	if err != nil {
		h.logger.Printf("ERROR: GetCountdownByID: %v", err)
		utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	fmt.Println(countdownData)

}