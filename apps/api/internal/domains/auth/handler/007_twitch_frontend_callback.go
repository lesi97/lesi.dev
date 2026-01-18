package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleTwitchFrontendAuthCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")

	if code == "" || state == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	stateCookie, err := r.Cookie("twitch_fe_state")
	if err != nil || stateCookie.Value == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	pkceCookie, err := r.Cookie("twitch_fe_pkce")
	if err != nil || pkceCookie.Value == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	identity, err := h.store.TwitchFrontendCallback(code, state, stateCookie.Value, pkceCookie.Value)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err)
		return
	}

	userID, err := h.store.TwitchFrontendUpsertUser(r.Context(), *identity)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	sessionToken, err := h.store.TwitchFrontendCreateSession(r.Context(), userID, 365*24*time.Hour)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	secure := os.Getenv("GO_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "lesidev_session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 365,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{Name: "twitch_fe_state", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode})
	http.SetCookie(w, &http.Cookie{Name: "twitch_fe_pkce", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode})

	http.Redirect(w, r, "/aim-trainer", http.StatusFound)
}