package plex

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
)

func (s *Store) findSeasonRatingKey(ctx context.Context, grandparentRatingKey string, season int) (*string, error) {
	url, err := s.appendXToken(fmt.Sprintf("/library/metadata/%s/children", grandparentRatingKey))
	if err != nil {
		return nil, err
	}

	raw, err := s.plexGET(ctx, *url)
	if err != nil {
		return nil, err
	}

	var res model.PlexLabels
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}

	for _, m := range res.MediaContainer.Metadata {
		if m.Index == season {
			key := m.RatingKey
			return &key, nil
		}
	}

	return nil, fmt.Errorf("season %d not found for show %s", season, grandparentRatingKey)
}
