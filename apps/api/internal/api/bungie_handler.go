package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/store/bungie_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type BungieHandler struct {
	logger         *utils.Logger
	bungieStore 	bungie_store.BungieStoreInterface
}

const bungieContextKey = "bungie"

func NewBungieHandler(logger *utils.Logger, bungieStore bungie_store.BungieStoreInterface)  *BungieHandler {
	return &BungieHandler{
		logger: logger,
		bungieStore: bungieStore,
	}
}

func (h *BungieHandler) HandleGetPlayTime(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("ERROR: handleSearchUser: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}

	reqInfo := bungie_store.BungieContextInfo{
		Platform: platform,
		Gamertag: idParam,
		Handler:  "HandleGetPlayTime",
	}
	ctx := context.WithValue(r.Context(), bungieContextKey, reqInfo)

	message, err := h.bungieStore.GetCharacterPlayTime(ctx)
	if err != nil {
		h.logger.Printf("ERROR: getCharacterPlayTime: %v", err)
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
	}
	utils.TextResponse(w, http.StatusOK, *message)
}

func (h *BungieHandler) HandleGetPrimary(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("ERROR: handleSearchUser: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}

	reqInfo := bungie_store.BungieContextInfo{
		Platform: platform,
		Gamertag: idParam,
		Handler:  "HandleGetPrimary",
		WeaponIndex: 0,
	}
	ctx := context.WithValue(r.Context(), bungieContextKey, reqInfo)

	message, err := h.bungieStore.GetEquippedWeapon(ctx)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}

func (h *BungieHandler) HandleGetSecondary(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("ERROR: handleSearchUser: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}

	reqInfo := bungie_store.BungieContextInfo{
		Platform: platform,
		Gamertag: idParam,
		Handler:  "HandleGetSecondary",
		WeaponIndex: 1,
	}
	ctx := context.WithValue(r.Context(), bungieContextKey, reqInfo)

	message, err := h.bungieStore.GetEquippedWeapon(ctx)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}

func (h *BungieHandler) HandleGetHeavy(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		h.logger.Printf("ERROR: handleSearchUser: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}

	reqInfo := bungie_store.BungieContextInfo{
		Platform: platform,
		Gamertag: idParam,
		Handler:  "HandleGetHeavy",
		WeaponIndex: 2,
	}
	ctx := context.WithValue(r.Context(), bungieContextKey, reqInfo)


	message, err := h.bungieStore.GetEquippedWeapon(ctx)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}

func (h *BungieHandler) HandleGetTerrorKillCount(w http.ResponseWriter, r *http.Request) {
	weaponName := r.URL.Query().Get("name")
	if weaponName == "" {
		utils.TextResponse(w, http.StatusBadRequest, "you must specify a weapon name")
		return
	}

	reqInfo := bungie_store.BungieContextInfo{
		Handler:  "HandleGetTerrorKillCount",
		WeaponName: weaponName,
	}
	ctx := context.WithValue(r.Context(), bungieContextKey, reqInfo)

	message, err := h.bungieStore.GetTerrorWeapon(ctx)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}