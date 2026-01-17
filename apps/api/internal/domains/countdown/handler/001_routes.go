package handler

import "github.com/go-chi/chi/v5"

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Route("/countdown", func(r chi.Router) {
		r.Post("/", h.createCountdown)
		r.Get("/{id}", h.getCountdown)
	})
}