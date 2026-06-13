package handler

import (
	"context"
	"net/http"

	bungie_store "github.com/lesi97/lesi.dev/internal/domains/bungie/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetTerrorKillCount(w http.ResponseWriter, r *http.Request) {
	weaponName := r.URL.Query().Get("name")
	if weaponName == "" {
		utils.TextResponse(w, http.StatusBadRequest, "you must specify a weapon name")
		return
	}

	reqInfo := bungie_store.BungieContextInfo{
		Handler:    "HandleGetTerrorKillCount",
		WeaponName: weaponName,
	}
	ctx := context.WithValue(r.Context(), bungieContextKey, reqInfo)

	message, err := h.store.GetTerrorWeapon(ctx)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	textResponse(w, http.StatusOK, *message)
}
