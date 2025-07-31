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
	h.bungieStore.GetEquippedPrimary(r.Context(), idParam, platform)
	// h.bungieStore.SearchDestinyPlayer(r.Context(), "Lesi%235934")
	// h.bungieStore.SearchDestinyPlayer(r.Context(), "4611686018475555326")

	utils.TextResponse(w, http.StatusOK, "All good")
}