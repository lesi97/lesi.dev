package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/model"
)


func (s *Store) UpsertUpdateUser(ctx context.Context, in model.UpdateInput) (*model.AimTrainerUpdate, error) {
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

	out := &model.AimTrainerUpdate{
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