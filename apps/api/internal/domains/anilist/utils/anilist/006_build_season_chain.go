package anilist

import (
	"context"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
)

func (s *Store) BuildSeasonChain(ctx context.Context, mediaID int) ([]model.AnilistRelationNode, error) {
	seen := make(map[int]bool)
	currentID := mediaID

	for {
		if seen[currentID] {
			break
		}
		seen[currentID] = true

		m, err := s.GetMediaRelationsByID(ctx, currentID)
		if err != nil {
			return nil, err
		}

		prequel, ok := pickPrequel(m.Relations.Edges)
		if !ok {
			break
		}
		if prequel.ID == 0 {
			break
		}
		currentID = prequel.ID
	}

	rootID := currentID

	chain := make([]model.AnilistRelationNode, 0)
	seen = make(map[int]bool)
	currentID = rootID

	for {
		if seen[currentID] {
			break
		}
		seen[currentID] = true

		m, err := s.GetMediaRelationsByID(ctx, currentID)
		if err != nil {
			return nil, err
		}

		chain = append(chain, model.AnilistRelationNode{
			ID:         m.ID,
			Title:      m.Title,
			SeasonYear: m.SeasonYear,
			Episodes:   m.Episodes,
			Format:     m.Format,
			IsAdult:    m.IsAdult,
		})

		sequels := pickSequels(m.Relations.Edges)
		if len(sequels) == 0 {
			break
		}
		if sequels[0].ID == 0 {
			break
		}
		currentID = sequels[0].ID
	}

	if len(chain) == 0 {
		return nil, fmt.Errorf("empty season chain")
	}

	return chain, nil
}
