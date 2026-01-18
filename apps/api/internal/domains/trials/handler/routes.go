package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/d2/trials", func(r chi.Router) {
		r.Get("/loot", h.HandleGetLoot)
		r.Get("/playercount", h.HandleGetPlayerCount)
	})
}
