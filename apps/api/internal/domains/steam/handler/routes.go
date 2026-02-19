package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/steam", func(r chi.Router) {
		r.Get("/playercount", h.HandleGetPlayerCount)
	})
}
