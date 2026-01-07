package api

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type TrialsHandler struct {
	logger     *utils.Logger
	TrialsStore store.TrialsStore
}

func NewTrialsHandler(logger *utils.Logger, TrialsStore store.TrialsStore)  *TrialsHandler {
	return &TrialsHandler{
		logger: logger,
		TrialsStore: TrialsStore,
	}
}

func (h *TrialsHandler) HandleGetLoot(w http.ResponseWriter, r *http.Request) {
	message := h.TrialsStore.GetLoot()
	utils.TextResponse(w, http.StatusOK, *message)
}

func (h *TrialsHandler) HandleGetPlayercount(w http.ResponseWriter, r *http.Request) {
	message := h.TrialsStore.GetPlayerCount()
	utils.TextResponse(w, http.StatusOK, *message)
}