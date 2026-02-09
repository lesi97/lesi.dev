package middleware

import (
	"net/http"
	"time"

	api_logs_store "github.com/lesi97/lesi.dev/internal/domains/api_logs/store"
	requestmetrics "github.com/lesi97/lesi.dev/internal/request_metrics"
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
	fetchCallsDuration := requestmetrics.GetFetchCallsDuration(r.Context())
	databaseCallsDuration := requestmetrics.GetDatabaseCallsDuration(r.Context())
	apiProcessingDuration := duration - fetchCallsDuration - databaseCallsDuration
	if apiProcessingDuration < 0 {
		apiProcessingDuration = 0
	}

	nonceElapsed, nonceOk := GetNonceElapsed(r.URL.Query().Get("nonce"), now)
	logLine := GetLogLine(path, sw.status, duration, nonceElapsed, nonceOk)
	responseBody := GetResponseBody(sw.body)
	hasNightbotHeaders := r.Header.Get("Nightbot-User") != "" || r.Header.Get("Nightbot-Channel") != ""
	hasStreamElementsHeader := r.Header.Get("X-Streamelements-Channel") != ""

	LogNightbotDetails(logger, r.Header, responseBody, sw.status)
	LogStreamElementsChannel(logger, r.Header, responseBody)

	if apiLogStore != nil && path != "/v1/telemetry" {
		LogApiRequest(
			apiLogStore,
			r,
			path,
			responseBody,
			duration,
			apiProcessingDuration,
			fetchCallsDuration,
			databaseCallsDuration,
			nonceElapsed,
			nonceOk,
			sw.status,
		)
	}
	if sw.status == http.StatusNotFound {
		LogNotFoundRequest(logger, logLine, r.Header, hasNightbotHeaders, hasStreamElementsHeader)
		return
	}

	logger.Printf("%v", logLine)
}
