package handler

import (
	"net/http"
	"os"

	"github.com/lesi97/lesi.dev/internal/ui/html"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleSpotifyAuthInitialRedirect(w http.ResponseWriter, r *http.Request) {
	if !isLocalAuthRequest(r) {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	secret := r.URL.Query().Get("secret")
	if secret == "" || secret != os.Getenv("PLEX_WEBHOOK_SECRET") {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	url, err := h.store.SpotifyAuthUrl(r.Context())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := html.GenerateAuthButton(*url, "Authenticate with Spotify")
	_, _ = w.Write([]byte(html))
}
