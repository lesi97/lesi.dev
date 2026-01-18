package anilist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
)

func (s *Store) SearchTitle(ctx context.Context, title string) ([]model.AnilistMedia, error) {
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
      isAdult
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

	raw, err := s.anilistPOST(
		ctx,
		s.env.GraphqlUrl,
		query,
		map[string]interface{}{
			"title": title,
		},
	)
	if err != nil {
		return nil, err
	}

	var res model.AnilistSearchResponse
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("anilist graphql error %s", res.Errors[0].Message)
	}

	return res.Data.Page.Media, nil
}
