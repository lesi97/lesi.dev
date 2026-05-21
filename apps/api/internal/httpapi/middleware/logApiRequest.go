package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/api_logs/model"
	api_logs_store "github.com/lesi97/lesi.dev/internal/domains/api_logs/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func LogApiRequest(
	apiLogStore api_logs_store.Methods,
	r *http.Request,
	path string,
	responseBody string,
	duration time.Duration,
	apiProcessingDuration time.Duration,
	fetchCallsDuration time.Duration,
	databaseCallsDuration time.Duration,
	nonceElapsed time.Duration,
	nonceOk bool,
	status int,
) {
	headers := r.Header.Clone()
	query := r.URL.Query()
	ip := utils.GetRequestIP(r)
	response := responseBody
	utils.TruncateString(&response, 2000)

	var nonceElapsedMS *int64
	if nonceOk {
		value := nonceElapsed.Milliseconds()
		nonceElapsedMS = &value
	}

	logEntry := model.ApiLog{
		Timestamp:       time.Now().UTC(),
		Route:           path,
		IP:              ip,
		Response:        response,
		ExecutionTimeMS: duration.Milliseconds(),
		ApiProcessingMS: apiProcessingDuration.Milliseconds(),
		FetchCallsMS:    fetchCallsDuration.Milliseconds(),
		DatabaseCallsMS: databaseCallsDuration.Milliseconds(),
		NonceElapsedMS:  nonceElapsedMS,
		StatusCode:      &status,
	}

	go func() {
		botType, channel, user := GetBotMetadata(headers, query)
		logEntry.BotType = botType
		logEntry.Channel = channel
		logEntry.User = user

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = apiLogStore.InsertApiLog(ctx, logEntry)
	}()
}
