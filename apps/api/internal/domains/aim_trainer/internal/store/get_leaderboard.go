package store

import (
	"context"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/internal/model"
)



func (s *Store) GetLeaderboard(ctx context.Context) ([]model.LeaderboardRow, error) {
	defer s.Logger.LogExecutionTime("DATABASE CALL: getLeaderboard", time.Now(), nil)
	const q = `
		SELECT username, score, accuracy
		FROM aim_trainer
		ORDER BY score DESC
	`

	rows, err := s.DB.Pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.LeaderboardRow, 0, 50)

	for rows.Next() {
		var r model.LeaderboardRow

		err = rows.Scan(
			&r.Username,
			&r.Score,
			&r.Accuracy,
		)
		if err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return out, nil
}
