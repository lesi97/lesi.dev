package anilist

import (
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

func pickPrequel(edges []model.AnilistRelationEdgeNode) (model.AnilistRelationNode, bool) {
	for _, e := range edges {
		if strings.EqualFold(e.RelationType, "PREQUEL") {
			return e.Node, true
		}
	}
	return model.AnilistRelationNode{}, false
}
