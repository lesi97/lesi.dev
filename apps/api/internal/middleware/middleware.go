package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Measure(logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.RequestURI()

			sw := &statusResponseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			defer func() {
				pathColour := utils.Colours["brightBlue"]
				timeColour := utils.Colours["green"]
				statusColour := utils.Colours["green"]

				duration := time.Since(start)

				if duration > 100*time.Millisecond {
					timeColour = utils.Colours["brightRed"] + utils.Colours["bold"]
				}

				if sw.status >= 400 {
					statusColour = utils.Colours["brightRed"] + utils.Colours["bold"]
				}

				statusBlock := fmt.Sprintf("%vStatus: %v%v", statusColour, sw.status, utils.Colours["reset"])
				pathBlock := fmt.Sprintf("%v%v%v", pathColour, path, utils.Colours["reset"])
				timeBlock := fmt.Sprintf("%vtook %v%v", timeColour, duration, utils.Colours["reset"])
				logger.Printf("%v | %v %v", statusBlock, pathBlock, timeBlock)
			}()

			next.ServeHTTP(sw, r)
		})
	}
}
