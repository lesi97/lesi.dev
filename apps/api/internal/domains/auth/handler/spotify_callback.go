package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleSpotifyAuthCallback(w http.ResponseWriter, r *http.Request) {
	if !isLocalAuthRequest(r) {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := h.store.SpotifyCallback(r.Context(), code); err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusNoContent, "updated")
}
