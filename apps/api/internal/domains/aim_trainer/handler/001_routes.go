package handler

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/aim-trainer", func(r chi.Router) {
		r.Get("/", h.getLeaderboard)
		r.Post("/", h.updateLeaderboard)
	})
}
