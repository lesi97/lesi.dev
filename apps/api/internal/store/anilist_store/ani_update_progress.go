package anilist_store

import (
	"encoding/json"
	"fmt"
)

type saveEntryData struct {
	SaveMediaListEntry struct {
		ID       int `json:"id"`
		MediaID  int `json:"mediaId"`
		Progress int `json:"progress"`
	} `json:"SaveMediaListEntry"`
}

type saveEntryResponse struct {
	Data   saveEntryData  `json:"data"`
	Errors []graphqlError `json:"errors,omitempty"`
}

func (s *AnilistStore) updateAniListProgress(mediaID int, progress int) error {
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
		s.graphql_url,
		query,
		map[string]interface{}{
			"mediaId":  mediaID,
			"progress": progress,
		},
	)
	if err != nil {
		return err
	}

	var res saveEntryResponse
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return err
	}

	if len(res.Errors) > 0 {
		return fmt.Errorf("anilist graphql error %s", res.Errors[0].Message)
	}

	return nil
}
