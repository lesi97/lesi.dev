package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/model"
	domain_utils "github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/utils"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) PostDateMemeClick(w http.ResponseWriter, r *http.Request) {
	if err := domain_utils.ValidateTelemetryAccess(r); err != nil {
		h.logger.Printf("Date meme click telemetry access denied: %v", err)
		utils.TextResponse(w, http.StatusForbidden, "telemetry access denied")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var payload model.DateMemeClickPayload
	if err := dec.Decode(&payload); err != nil {
		utils.TextResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := domain_utils.ValidateDateMemeClickPayload(payload); err != nil {
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.UpsertDateMemeClick(r.Context(), model.DateMemeClickInput{
		Route:     payload.Route,
		Action:    payload.Action,
		IP:        utils.GetRequestIP(r),
		UserAgent: r.UserAgent(),
	}); err != nil {
		h.logger.Printf("ERROR: Date meme click telemetry POST: %v", err)
		utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
