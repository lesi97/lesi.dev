package anilist

import (
	"sort"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

func pickSequels(edges []model.AnilistRelationEdgeNode) []model.AnilistRelationNode {
	out := make([]model.AnilistRelationNode, 0)
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
