package aim_trainer_store

import "context"

func (s *AimTrainerStore) upsert(ctx context.Context, u *AimTrainerUpdate) error {
	const q = `
		INSERT INTO aim_trainer (
			username,
			updated_at,
			completed_rounds,
			score,
			accuracy,
			milo_attacks,
			nice_triggers
		)
		VALUES (
			$1,
			$2,
			$3,
			COALESCE($4, 0),
			COALESCE($5, 0),
			COALESCE($6, 0),
			COALESCE($7, 0)
		)
		ON CONFLICT (username) DO UPDATE SET
			updated_at = EXCLUDED.updated_at,
			completed_rounds = EXCLUDED.completed_rounds,
			score = CASE
				WHEN $4 IS NULL THEN aim_trainer.score
				ELSE EXCLUDED.score
			END,
			accuracy = CASE
				WHEN $5 IS NULL THEN aim_trainer.accuracy
				ELSE EXCLUDED.accuracy
			END,
			milo_attacks = CASE
				WHEN $6 IS NULL THEN aim_trainer.milo_attacks
				ELSE EXCLUDED.milo_attacks
			END,
			nice_triggers = CASE
				WHEN $7 IS NULL THEN aim_trainer.nice_triggers
				ELSE EXCLUDED.nice_triggers
			END
	`

	var sVal *int64 = u.Score
	var aVal *float64 = u.Accuracy
	var mVal *int64 = u.MiloAttacks
	var nVal *int64 = u.NiceTriggers

	_, err := s.DB.Exec(ctx, q,
		u.Username,
		u.UpdatedAt,
		u.CompletedRounds,
		sVal,
		aVal,
		mVal,
		nVal,
	)

	return err
}