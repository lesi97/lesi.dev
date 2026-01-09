package anilist_store

import (
	"encoding/json"
	"fmt"
)

type relationEdgeNode struct {
	RelationType string    `json:"relationType"`
	Node         nodesType `json:"node"`
}

type mediaRelations struct {
	ID         int       `json:"id"`
	Title      titleType `json:"title"`
	Format     string    `json:"format"`
	SeasonYear *int      `json:"seasonYear"`
	Episodes   *int      `json:"episodes"`
	IsAdult    bool      `json:"isAdult"`
	Relations  struct {
		Edges []relationEdgeNode `json:"edges"`
	} `json:"relations"`
}


type mediaRelationsData struct {
	Media mediaRelations `json:"Media"`
}

type mediaRelationsResponse struct {
	Data   mediaRelationsData `json:"data"`
	Errors []graphqlError     `json:"errors,omitempty"`
}

func (s *AnilistStore) fetchMediaRelationsByID(mediaID int) (*mediaRelations, error) {
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
		s.graphql_url,
		query,
		map[string]interface{}{
			"id": mediaID,
		},
	)
	if err != nil {
		return nil, err
	}

	var res mediaRelationsResponse
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("anilist graphql error %s", res.Errors[0].Message)
	}

	return &res.Data.Media, nil
}
