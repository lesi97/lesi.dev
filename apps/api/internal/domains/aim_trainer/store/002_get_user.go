package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type AimTrainerRow struct {
	Username        string
	Score           float64
	Accuracy        float64
	CompletedRounds int64
	MiloAttacks     int64
	NiceTriggers    int64
}

func (s *Store) GetUser(ctx context.Context, username string) (*AimTrainerRow, error) {
	const q = `
		SELECT username, score, accuracy, completed_rounds, milo_attacks, nice_triggers
		FROM aim_trainer
		WHERE username = $1
		LIMIT 1
	`

	var row AimTrainerRow
	err := s.DB.QueryRow(ctx, q, username).Scan(
		&row.Username,
		&row.Score,
		&row.Accuracy,
		&row.CompletedRounds,
		&row.MiloAttacks,
		&row.NiceTriggers,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}