package middleware

import (
	"net/http"
	"time"

	api_logs_store "github.com/lesi97/lesi.dev/internal/domains/api_logs/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func LogRequestSummary(
	logger *utils.Logger,
	apiLogStore api_logs_store.Methods,
	sw *statusResponseWriter,
	r *http.Request,
	start time.Time,
	path string,
) {
	if path == "/favicon.ico" {
		return
	}

	now := time.Now()
	duration := now.Sub(start)

	nonceElapsed, nonceOk := GetNonceElapsed(r.URL.Query().Get("nonce"), now)
	logLine := GetLogLine(path, sw.status, duration, nonceElapsed, nonceOk)
	responseBody := GetResponseBody(sw.body)
	hasNightbotHeaders := r.Header.Get("Nightbot-User") != "" || r.Header.Get("Nightbot-Channel") != ""
	hasStreamElementsHeader := r.Header.Get("X-Streamelements-Channel") != ""

	LogNightbotDetails(logger, r.Header, responseBody, sw.status)
	LogStreamElementsChannel(logger, r.Header, responseBody)

	if apiLogStore != nil && path != "/v1/telemetry" {
		LogApiRequest(apiLogStore, r, path, responseBody, duration, nonceElapsed, nonceOk)
	}
	if sw.status == http.StatusNotFound {
		LogNotFoundRequest(logger, logLine, r.Header, hasNightbotHeaders, hasStreamElementsHeader)
		return
	}

	logger.Printf("%v", logLine)
}
