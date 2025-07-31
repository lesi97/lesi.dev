package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/lesi97/api.lesi.dev/internal/utils"
)

func Measure(logger *log.Logger, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		defer func() {
			pathColour := utils.Colours["brightBlack"]
			timeColour := utils.Colours["green"]
			duration := time.Since(start)
			if duration > 100*time.Millisecond {
				timeColour = utils.Colours["brightRed"] + utils.Colours["bold"]
			}
			logger.Printf("%s%s %stook %s%s", pathColour, path, timeColour, duration, utils.Colours["reset"])
		}()
		handlerFunc(w, r)
	}
}
