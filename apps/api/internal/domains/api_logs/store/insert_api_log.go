package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/domains/api_logs/model"
)

func (s *Store) InsertApiLog(ctx context.Context, log model.ApiLog) error {
	query := `
		INSERT INTO logs.api (
			event_timestamp,
			route,
			ip,
			channel,
			user_name,
			bot_type,
			response,
			execution_time_ms,
			api_processing_time_ms,
			fetch_calls_time_ms,
			database_calls_time_ms,
			nonce_elapsed_ms,
			status_code
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := s.DB.Exec(
		ctx,
		query,
		log.Timestamp,
		log.Route,
		log.IP,
		log.Channel,
		log.User,
		log.BotType,
		log.Response,
		log.ExecutionTimeMS,
		log.ApiProcessingMS,
		log.FetchCallsMS,
		log.DatabaseCallsMS,
		log.NonceElapsedMS,
		log.StatusCode,
	)
	if err != nil {
		return err
	}

	return nil
}
