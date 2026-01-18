package handler

import (
	"net/http"
	"os"

	"github.com/lesi97/lesi.dev/internal/html_pages"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleAniAuthInitialRedirect(w http.ResponseWriter, r *http.Request) {
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