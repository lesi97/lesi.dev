package api

import (
	"net/http"
	"os"

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

	url, err := h.store.GenerateAnilistAuthUrl()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(*url))
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
