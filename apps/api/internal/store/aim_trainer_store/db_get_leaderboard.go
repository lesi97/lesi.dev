package aim_trainer_store

import "context"

type LeaderboardRow struct {
	Username string
	Score    int64
	Accuracy float64
}

func (s *AimTrainerStore) GetLeaderboard(ctx context.Context) ([]LeaderboardRow, error) {
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

	out := make([]LeaderboardRow, 0, 50)

	for rows.Next() {
		var r LeaderboardRow

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