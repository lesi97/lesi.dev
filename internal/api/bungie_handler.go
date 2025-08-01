package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/api.lesi.dev/internal/store/bungie_store"
	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type BungieHandler struct {
	logger         *log.Logger
	bungieStore bungie_store.BungieStore
}

func NewBungieHandler(logger *log.Logger, bungieStore bungie_store.BungieStore)  *BungieHandler {
	return &BungieHandler{
		logger: logger,
		bungieStore: bungieStore,
	}
}

func (h *BungieHandler) HandleGetPlayTime(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	bungieGT := chi.URLParam(r, "id")
	if bungieGT == "" {
		h.logger.Printf("ERROR: handleSearchUser: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "invalid ID")
		return
	}
	message, err := h.bungieStore.GetCharacterPlayTime(r.Context(), bungieGT, platform)
	if err != nil {
		h.logger.Printf("ERROR: getCharacterPlayTime: %v", err)
		utils.TextResponse(w, http.StatusInternalServerError, "internal server error")
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
	message, err := h.bungieStore.GetEquippedWeapon(r.Context(), idParam, platform, 0)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "internal server error")
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
	message, err := h.bungieStore.GetEquippedWeapon(r.Context(), idParam, platform, 1)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "internal server error")
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
	message, err := h.bungieStore.GetEquippedWeapon(r.Context(), idParam, platform, 2)
	if err != nil {
		h.logger.Printf("ERROR: GetEquippedWeapon: No ID provided")
		utils.TextResponse(w, http.StatusBadRequest, "internal server error")
		return
	}

	utils.TextResponse(w, http.StatusOK, *message)
}