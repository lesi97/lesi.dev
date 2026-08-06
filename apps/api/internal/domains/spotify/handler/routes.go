package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/scrobbles", func(r chi.Router) {
		r.Get("/latest", h.GetLatestScrobble)
		r.Post("/", h.PostScrobble)
		r.Post("/poll/spotify", h.PostSpotifyPoll)
	})
}
