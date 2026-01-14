package api

import (
	"net/http"
	"os"
	"time"

	"github.com/lesi97/lesi.dev/internal/html_pages"
	"github.com/lesi97/lesi.dev/internal/store/auth_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AuthHandler struct {
	logger         *utils.Logger
	store 			auth_store.AuthStoreInterface
}

func NewAuthHandler(logger *utils.Logger, store auth_store.AuthStoreInterface) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		store:  store,
	}
}

func (h *AuthHandler) HandleAniAuthInitialRedirect(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	if secret == "" || secret != os.Getenv("PLEX_WEBHOOK_SECRET") {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	url, err := h.store.AnilistAuthUrl()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := html_pages.GenerateAuthButton(*url, "Authenticate with Anilist")
	_, _ = w.Write([]byte(html))
}

func (h *AuthHandler) HandleAnilistAuthCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
    code := query.Get("code")
	if code == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
	}

	err :=h.store.AnilistCallback(code)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusNoContent, "updated")
}


func (h *AuthHandler) HandleTwitchModAuthInitialRedirect(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	if secret == "" || secret != os.Getenv("PLEX_WEBHOOK_SECRET") {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	url, err := h.store.TwitchModAuthUrl()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := html_pages.GenerateAuthButton(*url, "Authenticate with Twitch")

	_, _ = w.Write([]byte(html))
}

func (h *AuthHandler) HandleTwitchModAuthCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
    code := query.Get("code")
	if code == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
	}

	err :=h.store.TwitchModCallback(code)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusNoContent, "updated")

}

func (h *AuthHandler) HandleTwitchFrontendAuthStart(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) HandleTwitchFrontendAuthCallback(w http.ResponseWriter, r *http.Request) {
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


func (h *AuthHandler) HandleTwitchAuthMe(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) HandleTwitchAuthLogout(w http.ResponseWriter, r *http.Request) {
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