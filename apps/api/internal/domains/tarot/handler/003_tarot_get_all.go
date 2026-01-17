package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/tarot/model"
)


func (h *Handler) getAllCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	cards := h.store.GetAllCards()
	response := model.CardsResponse{
		Cards: cards,
	}
	json.NewEncoder(w).Encode(response)
}