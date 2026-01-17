package anilist_store

import (
	"fmt"
	"sort"
	"strings"
)

func pickPrequel(edges []relationEdgeNode) (nodesType, bool) {
	for _, e := range edges {
		if strings.EqualFold(e.RelationType, "PREQUEL") {
			return e.Node, true
		}
	}
	return nodesType{}, false
}

func pickSequels(edges []relationEdgeNode) []nodesType {
	out := make([]nodesType, 0)
	for _, e := range edges {
		if strings.EqualFold(e.RelationType, "SEQUEL") {
			out = append(out, e.Node)
		}
	}

	sort.Slice(out, func(i int, j int) bool {
		yi := 0
		yj := 0
		if out[i].SeasonYear != nil {
			yi = *out[i].SeasonYear
		}
		if out[j].SeasonYear != nil {
			yj = *out[j].SeasonYear
		}
		if yi != yj {
			return yi < yj
		}
		return out[i].ID < out[j].ID
	})

	return out
}

func (s *AnilistStore) buildSeasonChain(mediaID int) ([]nodesType, error) {
	seen := make(map[int]bool)
	currentID := mediaID

	for {
		if seen[currentID] {
			break
		}
		seen[currentID] = true

		m, err := s.fetchMediaRelationsByID(currentID)
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

	chain := make([]nodesType, 0)
	seen = make(map[int]bool)
	currentID = rootID

	for {
		if seen[currentID] {
			break
		}
		seen[currentID] = true

		m, err := s.fetchMediaRelationsByID(currentID)
		if err != nil {
			return nil, err
		}

		chain = append(chain, nodesType{
			ID:         m.ID,
			Title:      m.Title,
			SeasonYear: m.SeasonYear,
			Episodes:   m.Episodes,
			Format:     m.Format,
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
