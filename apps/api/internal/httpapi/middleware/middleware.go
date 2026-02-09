package middleware

import (
	"net/http"
	"time"

	api_logs_store "github.com/lesi97/lesi.dev/internal/domains/api_logs/store"
	requestmetrics "github.com/lesi97/lesi.dev/internal/request_metrics"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.ResponseWriter.Write(data)
}

func Measure(logger *utils.Logger, apiLogStore api_logs_store.Methods) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthcheck" {
				next.ServeHTTP(w, r)
				return
			}
			r = r.WithContext(requestmetrics.WithRequestMetrics(r.Context()))
			start := time.Now()
			path := r.URL.RequestURI()

			sw := &statusResponseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			defer func() {
				LogRequestSummary(logger, apiLogStore, sw, r, start, path)
			}()

			next.ServeHTTP(sw, r)
		})
	}
}
