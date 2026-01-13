package aim_trainer_store

import (
	"context"
)

type UpdateInput struct {
	Username        string
	Score           *float64
	Accuracy        *float64
	HasNiceTriggers bool
	HasMiloAttacks  bool
	UpdatedAtISO    string
}

type AimTrainerUpdate struct {
	Username        string   `json:"username"`
	UpdatedAt       string   `json:"updated_at"`
	CompletedRounds int64    `json:"completed_rounds"`
	Score           *float64   `json:"score,omitempty"`
	Accuracy        *float64 `json:"accuracy,omitempty"`
	NiceTriggers    *int64   `json:"nice_triggers,omitempty"`
	MiloAttacks     *int64   `json:"milo_attacks,omitempty"`
}

type ExistingRow struct {
	Score           float64
	Accuracy        float64
	CompletedRounds int64
	MiloAttacks     int64
	NiceTriggers    int64
}

func (s *AimTrainerStore) UpsertUpdateUser(ctx context.Context, in UpdateInput) (*AimTrainerUpdate, error) {
	existing, err := s.GetUser(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	rounds := int64(1)
	if existing != nil {
		rounds = existing.CompletedRounds + 1
	}

	var niceTriggers *int64
	if in.HasNiceTriggers {
		v := int64(1)
		if existing != nil {
			v = existing.NiceTriggers + 1
		}
		niceTriggers = &v
	}

	var miloAttacks *int64
	if in.HasMiloAttacks {
		v := int64(1)
		if existing != nil {
			v = existing.MiloAttacks + 1
		}
		miloAttacks = &v
	}

	score := in.Score
	accuracy := in.Accuracy

	if existing != nil && score != nil {
		if accuracy != nil && *score == existing.Score && *accuracy < existing.Accuracy {
			accuracy = nil
		}

		if *score <= existing.Score {
			score = nil
			accuracy = nil
		}
	}

	out := &AimTrainerUpdate{
		Username:        in.Username,
		UpdatedAt:       in.UpdatedAtISO,
		CompletedRounds: rounds,
		Score:           score,
		Accuracy:        accuracy,
		NiceTriggers:    niceTriggers,
		MiloAttacks:     miloAttacks,
	}

	if err := s.upsert(ctx, out); err != nil {
		return nil, err
	}

	return out, nil
}