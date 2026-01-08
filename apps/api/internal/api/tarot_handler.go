package api

import (
	"encoding/json"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/store/tarot_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type TarotHandler struct {
	logger 		*utils.Logger
	tarotStore 	tarot_store.TarotStoreInterface
}

func NewTarotHandler(logger *utils.Logger, tarotStore tarot_store.TarotStoreInterface)  *TarotHandler {
	return &TarotHandler{
		logger: logger,
		tarotStore: tarotStore,
	}
}

func (h *TarotHandler) HandleGetRandomTarot(w http.ResponseWriter, r *http.Request) {
	tarot := h.tarotStore.GetRandomTarot()
	utils.TextResponse(w, http.StatusOK, *tarot)
}

func (h *TarotHandler) HandleGetAllCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	cards := h.tarotStore.GetAllCards()
	response := struct{
		Cards *[]tarot_store.Card `json:"cards"`
	}{
		Cards: cards,
	}
	json.NewEncoder(w).Encode(response)
}