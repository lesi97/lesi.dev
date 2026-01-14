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
				if path == "/favicon.ico" {
					return		
				}
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

				reset := utils.Colours["reset"]

				statusBlock := fmt.Sprintf("%vStatus: %v%v", statusColour, sw.status, reset)
				pathBlock := fmt.Sprintf("%v%v%v", pathColour, path, reset)
				timeBlock := fmt.Sprintf("%vtook %v%v", timeColour, duration, reset)
				if sw.status == http.StatusNotFound {
					fmt.Println()
					logger.Printf("%v | %v %v", statusBlock, pathBlock, timeBlock)
					for key, values := range r.Header {
						for _, value := range values {
							keyBlock := fmt.Sprintf("%v%v%v: %v", utils.Colours["brightBlue"], utils.Colours["dim"], key, reset)
							logger.Printf("%v%v\n", keyBlock, value)
						}
					}
					fmt.Println()
				} else {
					logger.Printf("%v | %v %v", statusBlock, pathBlock, timeBlock)
				}
			}()

			next.ServeHTTP(sw, r)
		})
	}
}
