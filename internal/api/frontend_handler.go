package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FrontendHandler struct {
	logger     *log.Logger
}

func NewFrontendHandler(logger *log.Logger)  *FrontendHandler {
	return &FrontendHandler{
		logger: logger,
	}
}

func (h *FrontendHandler) HandleFrontend(staticDir http.Dir) http.HandlerFunc {
	fileServer := http.FileServer(staticDir)

	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(string(staticDir), r.URL.Path)
		stat, err := os.Stat(path)

		if err == nil && !stat.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(string(staticDir), "index.html"))
	}
}

func (h *FrontendHandler) HandlePublic(staticDir http.Dir) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(string(staticDir), r.URL.Path)
		switch {
		case strings.HasSuffix(path, ".js.br"):
			h.setAndServeCompressed(w, r, path, "application/javascript")
		case strings.HasSuffix(path, ".css.br"):
			h.setAndServeCompressed(w, r, path, "text/css")
			w.Header().Set("Content-Encoding", "br")
		case strings.HasSuffix(path, ".wasm.br"):
			h.setAndServeCompressed(w, r, path, "application/wasm")
		case strings.HasSuffix(path, ".data.br"), strings.HasSuffix(path, ".br"):
			h.setAndServeCompressed(w, r, path, "application/octet-stream")
		default:
			http.ServeFile(w, r, path)
		}
	}
}

func (h *FrontendHandler) setAndServeCompressed(w http.ResponseWriter, r *http.Request, path string, contentType string) {
	file, err := os.Open(path)
	if err != nil {
		http.NotFound(w, nil)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Encoding", "br")
	http.ServeContent(w, r, path, time.Now(), file)
}