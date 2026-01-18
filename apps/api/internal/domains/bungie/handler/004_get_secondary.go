package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	bungie_store "github.com/lesi97/lesi.dev/internal/domains/bungie/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetSecondary(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("ERROR: handleSearchUser: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}

	reqInfo := bungie_store.BungieContextInfo{
		Platform:    platform,
		Gamertag:    idParam,
		Handler:     "HandleGetSecondary",
		WeaponIndex: 1,
	}
	ctx := context.WithValue(r.Context(), bungieContextKey, reqInfo)

	message, err := h.store.GetEquippedWeapon(ctx)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}
