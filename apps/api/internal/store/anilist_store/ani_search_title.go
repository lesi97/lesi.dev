package anilist_store

import (
	"encoding/json"
)

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphqlError struct {
	Message string `json:"message"`
}

type graphqlResponse[T any] struct {
	Data   T              `json:"data"`
	Errors []graphqlError `json:"errors,omitempty"`
}

type titleType struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type edgesType struct {
	RelationType string `json:"relationType"`
}

type nodesType struct {
	ID         int       `json:"id"`
	Title      titleType `json:"title"`
	SeasonYear *int      `json:"seasonYear"`
	Episodes   *int      `json:"episodes"`
	Format     string    `json:"format"`
}

type relationsType struct {
	Edges []edgesType `json:"edges"`
	Nodes []nodesType `json:"nodes"`
}

type mediaData struct {
	ID         int           `json:"id"`
	Title      titleType     `json:"title"`
	Format     string        `json:"format"`
	SeasonYear *int          `json:"seasonYear"`
	Episodes   *int          `json:"episodes"`
	Status     string        `json:"status"`
	Relations  relationsType `json:"relations"`
}

type media struct {
	Media []mediaData `json:"media"`
}

type pageData struct {
	Page media `json:"Page"`
}

type searchTitleResponse struct { 
	Data pageData `json:"data`
}

func (s *AnilistStore) SearchAnilist(title string) ([]mediaData, error) {
	query := `
query ($title: String) {
  Page(perPage: 10) {
    media(search: $title, type: ANIME) {
      id
      title { romaji english native }
      format
      seasonYear
      episodes
      status
      relations {
        edges { relationType }
        nodes {
          id
          title { romaji english native }
          seasonYear
          episodes
          format
        }
      }
    }
  }
}
`

	raw, err := s.AnilistPOST(
		s.graphql_url,
		query,
		map[string]interface{}{
			"title": title,
		},
	)
	if err != nil {
		return nil, err
	}

	var res searchTitleResponse

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}

	return res.Data.Page.Media, nil
}

