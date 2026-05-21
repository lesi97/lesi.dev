package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleGetRandomFact(w http.ResponseWriter, r *http.Request) {
	fact, err := h.store.GetRandomFact(r.Context())
	h.writeTextFact(w, fact, err)
}

func (h *Handler) HandleGetRandomCheeseFact(w http.ResponseWriter, r *http.Request) {
	fact, err := h.store.GetRandomCheeseFact(r.Context())
	h.writeTextFact(w, fact, err)
}

func (h *Handler) writeTextFact(w http.ResponseWriter, fact *string, err error) {
	if err != nil {
		h.logger.Printf("ERROR: facts | %v", err.Error())
		utils.TextResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.TextResponse(w, http.StatusOK, *fact)
}
