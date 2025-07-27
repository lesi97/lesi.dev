package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/lesi97/api.lesi.dev/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {

	routes := chi.NewRouter()

	routes.Get("/tarot", app.TarotHandler.HandleGetRandomTarot)
	routes.Get("/countdown/{id}", app.CountdownHandler.HandleGetCountdown)

	return routes
}