package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/api.lesi.dev/internal/app"
	"github.com/lesi97/api.lesi.dev/internal/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {

	routes := chi.NewRouter()

	routes.Get("/tarot", http.HandlerFunc(
		middleware.Measure(app.Logger, app.TarotHandler.HandleGetRandomTarot),
	))
	routes.Get("/tarot/all", http.HandlerFunc(
		middleware.Measure(app.Logger, app.TarotHandler.HandleGetAllCards),
	))

	routes.Get("/countdown/{id}", http.HandlerFunc(
		middleware.Measure(app.Logger, app.CountdownHandler.HandleGetCountdown),
	))

	routes.Get("/time", http.HandlerFunc(
		middleware.Measure(app.Logger, app.TimeapiHandler.HandleGetDateTime),
	))

	routes.Get("/d2/{id}/time", http.HandlerFunc(
		middleware.Measure(app.Logger, app.BungieHandler.HandleGetPlayTime),
	))

	routes.Get("/d2/{id}/primary", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))
	routes.Get("/d2/{id}/kinetic", http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))

	return routes
}