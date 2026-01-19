package httpapi

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func AddScriptRoutes(r chi.Router) {
	dir := os.Getenv("SCRIPTS_DIR")
	if dir == "" {
		dir = "/scripts"
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}

	fs := http.FileServer(http.Dir(abs))

	r.Get("/scripts/*", func(w http.ResponseWriter, req *http.Request) {
		path := chi.URLParam(req, "*")
		if path == "" {
			http.NotFound(w, req)
			return
		}

		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.Header().Set("Content-Type", "text/plain")
		req.URL.Path = path
		fs.ServeHTTP(w, req)
	})
}
