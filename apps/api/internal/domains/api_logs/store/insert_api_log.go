package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/domains/api_logs/model"
)

func (s *Store) InsertApiLog(ctx context.Context, log model.ApiLog) error {
	query := `
		INSERT INTO logs.api (event_timestamp, route, ip, channel, user_name, bot_type, response, execution_time_ms, nonce_elapsed_ms)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
		log.NonceElapsedMS,
	)
	if err != nil {
		return err
	}

	return nil
}
