package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/countdown", func(r chi.Router) {
		r.Post("/", h.CreateCountdown)
		r.Get("/{id}", h.getCountdown)
	})
}
