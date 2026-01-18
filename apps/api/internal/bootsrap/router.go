package bootstrap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/httpapi/middleware"
)

type DisabledServiceHandler struct {
	ServiceName string
}

func (h *DisabledServiceHandler) RegisterRoutes(r chi.Router) {
	r.NotFound(h.HandleDisabled)
}

func (h *DisabledServiceHandler) HandleDisabled(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	_, _ = w.Write([]byte(h.ServiceName + " service is disabled"))
}

func handleIndex(r chi.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		_, _ = w.Write([]byte("Available routes\n\n"))
		_ = chi.Walk(r, func(method string, pattern string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
			_, _ = w.Write([]byte(method + " - " + pattern + "\n"))
			return nil
		})
	}
}

func setupRoutes(app *Application) *chi.Mux {

	routes := chi.NewRouter()
	routes.Use(middleware.Measure(app.Logger))
	routes.Route("/v1", func(r chi.Router) {
		app.tarot.RegisterRoutes(r)
		app.time.RegisterRoutes(r)
		if app.AnilistHandler != nil {
			app.AnilistHandler.RegisterRoutes(r)
		}
	})

	return routes
	// routes.Get("/healthcheck", api.Healthcheck)

	// routes.Post("/v1/countdown", http.HandlerFunc(app.CountdownHandler.HandleCountdownPost))
	// routes.Get("/v1/countdown/{id}", http.HandlerFunc(app.CountdownHandler.HandleGetCountdown))

	// routes.Get("/v1/d2/trials/loot", http.HandlerFunc(app.TrialsHandler.HandleGetLoot))
	// routes.Get("/v1/d2/playercount", http.HandlerFunc(app.TrialsHandler.HandleGetPlayercount))

	// routes.Get("/v1/aim-trainer", http.HandlerFunc(app.AimTrainerHandler.HandleGetLeaderboard))
	// routes.Post("/v1/aim-trainer", http.HandlerFunc(app.AimTrainerHandler.HandleUpdateLeaderboard))

	// if app.BungieHandler != nil {
	// 	routes.Get("/v1/d2/{id}/time", http.HandlerFunc(app.BungieHandler.HandleGetPlayTime))
	// 	routes.Get("/v1/d2/{id}/primary", http.HandlerFunc(app.BungieHandler.HandleGetPrimary))
	// 	routes.Get("/v1/d2/{id}/kinetic", http.HandlerFunc(app.BungieHandler.HandleGetPrimary))
	// 	routes.Get("/v1/d2/{id}/secondary", http.HandlerFunc(app.BungieHandler.HandleGetSecondary))
	// 	routes.Get("/v1/d2/{id}/energy", http.HandlerFunc(app.BungieHandler.HandleGetSecondary))
	// 	routes.Get("/v1/d2/{id}/heavy", http.HandlerFunc(app.BungieHandler.HandleGetHeavy))
	// 	routes.Get("/v1/d2/terror/weapons",	http.HandlerFunc(app.BungieHandler.HandleGetTerrorKillCount))
	// 	routes.Get("/v1/d2/yunger/weapons",	http.HandlerFunc(app.BungieHandler.HandleGetYungerKillCount))
	// 	routes.Get("/v1/d2/reset",	http.HandlerFunc(app.BungieHandler.HandleReset))
	// } else {
	// 	app.Logger.PrintColour(true, "brightRed", "Bungie routes not registered")
	// 	disabled := &DisabledServiceHandler{
	// 		ServiceName: "Bungie",
	// 	}
	// 	routes.Route("/v1/d2/terror", disabled.RegisterRoutes)
	// 	routes.Route("/v1/d2/{id}", disabled.RegisterRoutes)
	// }

	// if app.AnilistHandler != nil {
	// 	routes.Post("/v1/anilist", http.HandlerFunc(app.AnilistHandler.HandleUpdateAnilist))
	// } else {
	// 	app.Logger.PrintColour(true, "brightRed", "AniList routes not registered")
	// 	disabled := &DisabledServiceHandler{
	// 		ServiceName: "AniList",
	// 	}
	// 	if app.AuthHandler != nil {
	// 		routes.Get("/v1/auth/anilist/login", http.HandlerFunc(app.AuthHandler.HandleAniAuthInitialRedirect))
	// 		routes.Get("/v1/auth/anilist/callback", http.HandlerFunc(app.AuthHandler.HandleAnilistAuthCallback))
	// 	}

	// 	routes.Route("/v1/anilist", disabled.RegisterRoutes)
	// }

	// if app.TwitchHandler != nil {
	// 	routes.Get("/v1/twitch/{streamer}/chatters", http.HandlerFunc(app.TwitchHandler.HandleGetRandomChatter))
	// } else {
	// 	app.Logger.PrintColour(true, "brightRed", "Twitch routes not registered")
	// 	disabled := &DisabledServiceHandler{
	// 		ServiceName: "Twitch_GO",
	// 	}
	// 	// These routes are intended for the nightbot command using a mod user's twitch auth
	// 	if app.AuthHandler != nil {
	// 		routes.Get("/v1/auth/twitch/login", http.HandlerFunc(app.AuthHandler.HandleTwitchModAuthInitialRedirect))
	// 		routes.Get("/v1/auth/twitch/callback", http.HandlerFunc(app.AuthHandler.HandleTwitchModAuthCallback))
	// 	}

	// 	routes.Route("/v1/twitch", disabled.RegisterRoutes)
	// }
	// // Always allow the below routes to allow users to login to the aim trainer
	// routes.Get("/auth/twitch", app.AuthHandler.HandleTwitchFrontendAuthStart)
	// routes.Get("/auth/twitch/callback", app.AuthHandler.HandleTwitchFrontendAuthCallback)
	// routes.Get("/auth/twitch/me", app.AuthHandler.HandleTwitchAuthMe)
	// routes.Post("/auth/twitch/logout", app.AuthHandler.HandleTwitchAuthLogout)

	// if os.Getenv("GO_ENV") == "development" {
	// 	routes.Post("/local/db-dump", http.HandlerFunc(app.LocalHandler.HandleDbDump))
	// } else {
	// 	disabled := &DisabledServiceHandler{
	// 		ServiceName: "Local",
	// 	}
	// 	routes.Route("/local", disabled.RegisterRoutes)
	// }

	// if os.Getenv("GO_ENV") == "development" {
	// 	routes.Get("/", handleIndex(routes))
	// }

	return routes
}
