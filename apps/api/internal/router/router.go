package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/api/internal/app"
	"github.com/lesi97/lesi.dev/api/internal/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {

	routes := chi.NewRouter()

	routes.Get("/api/time", http.HandlerFunc(middleware.Measure(app.Logger, app.TimeapiHandler.HandleGetDateTime)))

	routes.Get("/api/tarot", 		http.HandlerFunc(middleware.Measure(app.Logger, app.TarotHandler.HandleGetRandomTarot)))
	routes.Get("/api/tarot/all", 	http.HandlerFunc(middleware.Measure(app.Logger, app.TarotHandler.HandleGetAllCards)))

	routes.Get("/api/countdown/{id}", http.HandlerFunc(middleware.Measure(app.Logger, app.CountdownHandler.HandleGetCountdown)))

	routes.Get("/api/d2/{id}/time", 		http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPlayTime)))
	routes.Get("/api/d2/{id}/primary", 		http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))
	routes.Get("/api/d2/{id}/kinetic", 		http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetPrimary)))
	routes.Get("/api/d2/{id}/secondary", 	http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetSecondary)))
	routes.Get("/api/d2/{id}/energy", 		http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetSecondary)))
	routes.Get("/api/d2/{id}/heavy", 		http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetHeavy)))
	routes.Get("/api/d2/terror/weapons", 	http.HandlerFunc(middleware.Measure(app.Logger, app.BungieHandler.HandleGetTerrorKillCount)))

	routes.Get("/api/d2/trials/loot", 		http.HandlerFunc(middleware.Measure(app.Logger, app.TrialsHandler.HandleGetLoot)))
	routes.Get("/api/d2/playercount", 		http.HandlerFunc(middleware.Measure(app.Logger, app.TrialsHandler.HandleGetPlayercount)))

	
	publicDir := http.Dir("./web/public")
	routes.Handle("/_static/*", app.FrontendHandler.HandlePublic(publicDir))
	
	distDir := http.Dir("./web/dist")
	routes.Handle("/*", app.FrontendHandler.HandleFrontend(distDir))

	return routes
}