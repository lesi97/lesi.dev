package store

import (
	"context"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/model"
)

func (s *Store) InsertTelemetry(ctx context.Context, payload model.TelemetryPayload) error {
	timestamp, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO logs.web (event_timestamp, route, user_agent, ip, error)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = s.DB.Exec(ctx, query, timestamp, payload.Route, payload.UserAgent, payload.IP, payload.Error)
	if err != nil {
		return err
	}

	return nil
}
