package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/spotify", func(r chi.Router) {
		r.Get("/latest", h.GetLatestScrobble)
		r.Post("/", h.PostScrobble)
		r.Post("/poll", h.PostSpotifyPoll)
		r.Post("/enrich", h.PostSpotifyEnrich)
	})
}
