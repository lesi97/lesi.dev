package handler

import (
	"net/http"
	"os"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleTwitchFrontendAuthStart(w http.ResponseWriter, r *http.Request) {
	res, err := h.store.TwitchFrontendAuthStart()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	secure := os.Getenv("GO_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "twitch_fe_state",
		Value:    res.State,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "twitch_fe_pkce",
		Value:    res.PKCEVerifier,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, res.URL, http.StatusFound)
}
