package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lesi97/api.lesi.dev/internal/store"
	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type TarotHandler struct {
	logger 		*log.Logger
	tarotStore 	store.TarotStore
}

func NewTarotHandler(logger *log.Logger, tarotStore store.TarotStore)  *TarotHandler {
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
		Cards *[]store.TarotCard `json:"cards"`
	}{
		Cards: cards,
	}
	json.NewEncoder(w).Encode(response)
}