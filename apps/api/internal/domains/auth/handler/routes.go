package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Get("/anilist/login", h.HandleAniAuthInitialRedirect)
		r.Get("/anilist/callback", h.HandleAnilistAuthCallback)
		r.Get("/spotify/login", h.HandleSpotifyAuthInitialRedirect)
		r.Get("/spotify/callback", h.HandleSpotifyAuthCallback)
		r.Get("/twitch/login", h.HandleTwitchModAuthInitialRedirect)
		r.Get("/twitch/callback", h.HandleTwitchModAuthCallback)
	})

	r.Route("/auth/twitch", func(r chi.Router) {
		r.Get("/", h.HandleTwitchFrontendAuthStart)
		r.Get("/callback", h.HandleTwitchFrontendAuthCallback)
		r.Get("/me", h.HandleTwitchAuthMe)
		r.Post("/logout", h.HandleTwitchAuthLogout)
	})
}
