package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/store/countdown_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type CountdownHandler struct {
	logger     *utils.Logger
	countdownStore countdown_store.CountdownStoreInterface
}

func NewCountdownHandler(logger *utils.Logger, countdownStore countdown_store.CountdownStoreInterface)  *CountdownHandler {
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

func (h *CountdownHandler) HandleCountdownPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req countdown_store.CountdownPostRequest
	if err := dec.Decode(&req); err != nil {
		utils.TextResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TargetDate.IsZero() {
		utils.TextResponse(w, http.StatusBadRequest, "target_date is required")
		return
	}

	if req.Message == "" {
		utils.TextResponse(w, http.StatusBadRequest, "message is required")
		return
	}

	if req.FallbackMessage == "" {
		utils.TextResponse(w, http.StatusBadRequest, "fallback_message is required")
		return
	}

	if req.TargetDate.Before(time.Now().UTC()) {
		utils.TextResponse(w, http.StatusBadRequest, "target_date must be in the future")
		return
	}

	uuid, err := h.countdownStore.InsertCountdown(r.Context(), req)
	if err != nil {
		h.logger.Printf("ERROR: Countdown POST: %v", err)
		utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.TextResponse(w, http.StatusOK, *uuid)
}