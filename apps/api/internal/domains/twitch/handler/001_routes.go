package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/twitch", func(r chi.Router) {
		r.Get("/{streamer}/chatters", h.HandleGetRandomChatter)
	})
}
