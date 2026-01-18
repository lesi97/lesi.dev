package handler

import (
	"net/http"
	"os"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleTwitchAuthLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("lesidev_session")
	if err == nil && c.Value != "" {
		_ = h.store.TwitchFrontendDeleteSessionByToken(r.Context(), c.Value)
	}

	secure := os.Getenv("GO_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "lesidev_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})

	utils.Success(w, http.StatusNoContent, "ok")
}