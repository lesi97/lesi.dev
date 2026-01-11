package api

import (
	"net/http"
	"os"

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


func (h *AuthHandler) HandleTwitchAuthInitialRedirect(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	if secret == "" || secret != os.Getenv("PLEX_WEBHOOK_SECRET") {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	url, err := h.store.TwitchAuthUrl()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := html_pages.GenerateAuthButton(*url, "Authenticate with Twitch")

	_, _ = w.Write([]byte(html))
}

func (h *AuthHandler) HandleTwitchlistAuthCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
    code := query.Get("code")
	if code == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
	}

	err :=h.store.TwitchCallback(code)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusNoContent, "updated")

}
