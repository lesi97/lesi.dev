package anilist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

type updateProgressResponse struct {
	Data struct {
		SaveMediaListEntry struct {
			ID       int `json:"id"`
			MediaID  int `json:"mediaId"`
			Progress int `json:"progress"`
		} `json:"SaveMediaListEntry"`
	} `json:"data"`
	Errors []model.AnilistGraphQLError `json:"errors,omitempty"`
}

func (s *Store) UpdateProgress(ctx context.Context, mediaID int, progress int) error {
	query := `
mutation ($mediaId: Int, $progress: Int) {
  SaveMediaListEntry(mediaId: $mediaId, progress: $progress) {
    id
    mediaId
    progress
  }
}
`

	raw, err := s.anilistPOST(
		ctx,
		s.env.GraphqlUrl,
		query,
		map[string]interface{}{
			"mediaId":  mediaID,
			"progress": progress,
		},
	)
	if err != nil {
		return err
	}

	var res updateProgressResponse
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return err
	}

	if len(res.Errors) > 0 {
		return fmt.Errorf("anilist graphql error %s", res.Errors[0].Message)
	}

	return nil
}
