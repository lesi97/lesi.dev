package anilist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
)

func (s *Store) GetMediaRelationsByID(ctx context.Context, mediaID int) (*model.AnilistMediaRelations, error) {
	query := `
query ($id: Int) {
  Media(id: $id, type: ANIME) {
    id
    title { romaji english native }
    format
    seasonYear
    episodes
    isAdult
    relations {
      edges {
        relationType
        node {
          id
          title { romaji english native }
          seasonYear
          episodes
          format
          isAdult
        }
      }
    }
  }
}
`

	raw, err := s.anilistPOST(
		ctx,
		s.env.GraphqlUrl,
		query,
		map[string]interface{}{
			"id": mediaID,
		},
	)
	if err != nil {
		return nil, err
	}

	var res model.AnilistMediaRelationsResponse
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("anilist graphql error %s", res.Errors[0].Message)
	}

	return &res.Data.Media, nil
}
