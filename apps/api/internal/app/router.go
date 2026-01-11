package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/middleware"
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

	routes.Get("/v1/time",	http.HandlerFunc(middleware.Measure(app.Logger, app.TimeapiHandler.HandleGetDateTime)))

	routes.Get("/v1/tarot", http.HandlerFunc(middleware.Measure(app.Logger, app.TarotHandler.HandleGetRandomTarot)))
	routes.Get("/v1/tarot/all", http.HandlerFunc(middleware.Measure(app.Logger, app.TarotHandler.HandleGetAllCards)))

	routes.Post("/v1/countdown", http.HandlerFunc(middleware.Measure(app.Logger, app.CountdownHandler.HandleCountdownPost)))
	routes.Get("/v1/countdown/{id}", http.HandlerFunc(middleware.Measure(app.Logger, app.CountdownHandler.HandleGetCountdown)))

	routes.Get("/v1/d2/trials/loot", http.HandlerFunc(middleware.Measure(app.Logger, app.TrialsHandler.HandleGetLoot)))
	routes.Get("/v1/d2/playercount", http.HandlerFunc(middleware.Measure(app.Logger, app.TrialsHandler.HandleGetPlayercount)))


	if app.BungieHandler != nil {
		routes.Get("/v1/d2/{id}/time", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPlayTime)))
		routes.Get("/v1/d2/{id}/primary", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))
		routes.Get("/v1/d2/{id}/kinetic", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))
		routes.Get("/v1/d2/{id}/secondary", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetSecondary)))
		routes.Get("/v1/d2/{id}/energy", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetSecondary)))
		routes.Get("/v1/d2/{id}/heavy", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetHeavy)))
		routes.Get("/v1/d2/terror/weapons",	http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetTerrorKillCount)))
		routes.Get("/v1/d2/yunger/weapons",	http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetYungerKillCount)))
		routes.Get("/v1/d2/reset",	http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleReset)))
	} else {
		app.Logger.PrintColour(true, "brightRed", "Bungie routes not registered")
		disabled := &DisabledServiceHandler{
			ServiceName: "Bungie",
		}
		routes.Route("/v1/d2/terror", disabled.RegisterRoutes)
		routes.Route("/v1/d2/{id}", disabled.RegisterRoutes)
	}


	if app.AnilistHandler != nil {
		routes.Post("/v1/anilist", http.HandlerFunc(middleware.Measure(app.Logger, app.AnilistHandler.HandleUpdateAnilist)))
	} else {
		app.Logger.PrintColour(true, "brightRed", "AniList routes not registered")
		disabled := &DisabledServiceHandler{
			ServiceName: "AniList",
		}
		routes.Get("/v1/auth/anilist/login", http.HandlerFunc(middleware.Measure(app.Logger, app.AuthHandler.HandleAniAuthInitialRedirect)))
		routes.Get("/v1/auth/anilist/callback", http.HandlerFunc(middleware.Measure(app.Logger, app.AuthHandler.HandleAnilistAuthCallback)))
		routes.Route("/v1/anilist", disabled.RegisterRoutes)
	}
	

	if app.TwitchHandler != nil {
		routes.Get("/v1/twitch/{streamer}/chatters", http.HandlerFunc(middleware.Measure(app.Logger, app.TwitchHandler.HandleGetRandomChatter)))
	} else {
		app.Logger.PrintColour(true, "brightRed", "Twitch routes not registered")
		disabled := &DisabledServiceHandler{
			ServiceName: "Twitch_GO",
		}
		routes.Get("/v1/auth/twitch/login", http.HandlerFunc(middleware.Measure(app.Logger, app.AuthHandler.HandleTwitchAuthInitialRedirect)))
		routes.Get("/v1/auth/twitch/callback", http.HandlerFunc(middleware.Measure(app.Logger, app.AuthHandler.HandleTwitchlistAuthCallback)))
		routes.Route("/v1/twitch", disabled.RegisterRoutes)
	}


	if os.Getenv("GO_ENV") == "development" {
		routes.Post("/local/db-dump", http.HandlerFunc(middleware.Measure(app.Logger, app.LocalHandler.HandleDbDump)))
	} else {
		disabled := &DisabledServiceHandler{
			ServiceName: "Local",
		}
		routes.Route("/local", disabled.RegisterRoutes)
	}

	routes.Get("/", handleIndex(routes))
	
	return routes
}