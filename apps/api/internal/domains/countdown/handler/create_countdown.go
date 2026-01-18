package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/countdown/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) CreateCountdown(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req model.PostData
	if err := dec.Decode(&req); err != nil {
		utils.TextResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	uuid, err := h.store.InsertCountdown(r.Context(), req)
	if err != nil {
		h.logger.Printf("ERROR: Countdown POST: %v", err)
		utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.TextResponse(w, http.StatusOK, *uuid)
}
