package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleTwitchAuthMe(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("lesidev_session")
	if err != nil || c.Value == "" {
		utils.Error(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	u, err := h.store.TwitchFrontendGetUserBySession(r.Context(), c.Value)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	utils.Success(w, http.StatusOK, u)
}