package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/d2", func(r chi.Router) {
		r.Get("/{id}/time", h.HandleGetPlayTime)
		r.Get("/{id}/primary", h.HandleGetPrimary)
		r.Get("/{id}/kinetic", h.HandleGetPrimary)
		r.Get("/{id}/secondary", h.HandleGetSecondary)
		r.Get("/{id}/energy", h.HandleGetSecondary)
		r.Get("/{id}/heavy", h.HandleGetHeavy)
		r.Get("/terror/weapons", h.HandleGetTerrorKillCount)
		r.Get("/yunger/weapons", h.HandleGetYungerKillCount)
		r.Get("/reset", h.HandleReset)
	})
}
