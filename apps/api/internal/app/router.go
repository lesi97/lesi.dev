package app

import (
	"net/http"

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

func setupRoutes(app *Application) *chi.Mux {

	routes := chi.NewRouter()

	routes.Get("/time",	http.HandlerFunc(middleware.Measure(app.Logger, app.TimeapiHandler.HandleGetDateTime)))

	routes.Get("/v1/tarot", http.HandlerFunc(middleware.Measure(app.Logger, app.TarotHandler.HandleGetRandomTarot)))
	routes.Get("/tarot/all", http.HandlerFunc(middleware.Measure(app.Logger, app.TarotHandler.HandleGetAllCards)))

	routes.Post("/v1/countdown", http.HandlerFunc(middleware.Measure(app.Logger, app.CountdownHandler.HandleCountdownPost)))
	routes.Get("/v1/countdown/{id}", http.HandlerFunc(middleware.Measure(app.Logger, app.CountdownHandler.HandleGetCountdown)))

	routes.Get("/v1/d2/{id}/time", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPlayTime)))
	routes.Get("/v1/d2/{id}/primary", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))
	routes.Get("/v1/d2/{id}/kinetic", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))
	routes.Get("/v1/d2/{id}/secondary", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetSecondary)))
	routes.Get("/v1/d2/{id}/energy", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetSecondary)))
	routes.Get("/v1/d2/{id}/heavy", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetHeavy)))
	routes.Get("/v1/d2/terror/weapons",	http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetTerrorKillCount)))

	routes.Get("/v1/d2/trials/loot", http.HandlerFunc(middleware.Measure(app.Logger, app.TrialsHandler.HandleGetLoot)))
	routes.Get("/v1/d2/playercount", http.HandlerFunc(middleware.Measure(app.Logger, app.TrialsHandler.HandleGetPlayercount)))

	if app.AnilistHandler != nil {
		routes.Get("/v1/anilist/test", http.HandlerFunc(middleware.Measure(app.Logger, app.AnilistHandler.HandleAnilistTest)))
	} else {
		app.Logger.PrintColour(true, "brightRed", "AniList routes not registered")
		disabled := &DisabledServiceHandler{
			ServiceName: "AniList",
		}
		routes.Route("/v1/anilist", disabled.RegisterRoutes)
	}
	
	return routes
}