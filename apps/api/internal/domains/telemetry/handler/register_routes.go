package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/telemetry", func(r chi.Router) {
		r.Post("/", h.PostTelemetry)
	})
}
