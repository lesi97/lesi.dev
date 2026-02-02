package bootstrap

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/httpapi"
	"github.com/lesi97/lesi.dev/internal/httpapi/middleware"
)

type DisabledServiceHandler struct {
	ServiceName string
}

type RouteRegistrar interface {
	RegisterRoutes(r chi.Router)
}

func registerRoutes(registrar RouteRegistrar, r chi.Router) {
	registrar.RegisterRoutes(r)
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
	routes.Use(middleware.RateLimit())
	
	routes.Get("/healthcheck", http.HandlerFunc(httpapi.Healthcheck))
	httpapi.AddScriptRoutes(routes)
	httpapi.AddFakeEnv(routes)

	routes.Route("/v1", func(r chi.Router) {
		if app.TarotHandler != nil {
			registerRoutes(app.TarotHandler, r)
		}
		if app.TimeHandler != nil {
			registerRoutes(app.TimeHandler, r)
		}
		if app.AnilistHandler != nil {
			registerRoutes(app.AnilistHandler, r)
		}
		if app.AimTrainerHandler != nil {
			registerRoutes(app.AimTrainerHandler, r)
		}
		if app.CountdownHandler != nil {
			registerRoutes(app.CountdownHandler, r)
		}
		if app.TrialsHandler != nil {
			registerRoutes(app.TrialsHandler, r)
		}
		if app.BungieHandler != nil {
			registerRoutes(app.BungieHandler, r)
		}
		if app.TwitchHandler != nil {
			registerRoutes(app.TwitchHandler, r)
		}

	})

	if app.AuthHandler != nil {
		registerRoutes(app.AuthHandler, routes)
	}

	if app.LocalHandler != nil && os.Getenv("GO_ENV") == "development" {
		registerRoutes(app.LocalHandler, routes)
	}

	if os.Getenv("GO_ENV") == "development" {
		routes.Get("/", handleIndex(routes))
	}

	return routes

}
