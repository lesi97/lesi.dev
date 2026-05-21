package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/facts", func(r chi.Router) {
		r.Get("/", h.HandleGetRandomFact)
		r.Get("/cheese", h.HandleGetRandomCheeseFact)
	})
}
