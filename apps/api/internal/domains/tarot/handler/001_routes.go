package handler

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router) {

	r.Route("/tarot", func(r chi.Router) {
		r.Get("/", h.getRandomTarot)
		r.Get("/all", h.getAllCards)
	})
}