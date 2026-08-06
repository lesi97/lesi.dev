package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	domain_utils "github.com/lesi97/lesi.dev/internal/domains/spotify/internal/utils"
	shared_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) PostScrobble(w http.ResponseWriter, r *http.Request) {
	if err := domain_utils.ValidateScrobbleAccess(r); err != nil {
		h.logger.Printf("Scrobble access denied: %v", err)
		shared_utils.Error(w, http.StatusForbidden, "scrobble access denied")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var raw json.RawMessage
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&raw); err != nil {
		shared_utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(raw) == 0 {
		shared_utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var req model.ScrobbleRequest
	dec = json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		shared_utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input, err := req.ToInput(raw)
	if err != nil {
		shared_utils.Error(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.store.InsertScrobble(r.Context(), input)
	if err != nil {
		h.logger.Printf("ERROR: Scrobble POST: %v", err)
		shared_utils.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if result == nil {
		h.logger.Print("ERROR: Scrobble POST: store returned nil result")
		shared_utils.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": result,
		"error":   nil,
	}); err != nil && !errors.Is(err, http.ErrHandlerTimeout) {
		h.logger.Printf("ERROR: Scrobble response encode: %v", err)
	}
}
