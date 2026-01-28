package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

func Measure(logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthcheck" {
				next.ServeHTTP(w, r)
				return
			}
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
				responseBody := strings.TrimSpace(string(sw.body))
				if responseBody == "" {
					responseBody = "<empty>"
				}
				hasNightbotHeaders := r.Header.Get("Nightbot-User") != "" || r.Header.Get("Nightbot-Channel") != ""
				hasStreamElementsHeader := r.Header.Get("X-Streamelements-Channel") != ""
				userDisplayName, channelDisplayName, ok := GetNightbotDisplayNames(r.Header)
				if ok {
					nightbotColour := utils.Colours["brightMagenta"]
					logger.Printf("%vNightbot User: %v%v", nightbotColour, userDisplayName, reset)
					logger.Printf("%vNightbot Channel: %v%v", nightbotColour, channelDisplayName, reset)
					logger.Printf("%vNightbot Response: %v%v", nightbotColour, responseBody, reset)
				}
				LogStreamElementsChannel(logger, r.Header, responseBody)
				if sw.status == http.StatusNotFound {
					if hasNightbotHeaders || hasStreamElementsHeader {
						return
					}
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
