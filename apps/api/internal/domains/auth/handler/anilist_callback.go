package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleAnilistAuthCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	if code == "" {
		utils.Error(w, http.StatusBadRequest, "bad request")
	}

	err := h.store.AnilistCallback(code)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusNoContent, "updated")
}
