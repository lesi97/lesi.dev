package plex

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
)

func (s *Store) getSeasonEpisodeCount(ctx context.Context, grandparentRatingKey string, season int) (int, error) {
	if _, ok := s.seasonCountCache[grandparentRatingKey]; !ok {
		s.seasonCountCache[grandparentRatingKey] = make(map[int]int)
	}
	if cached, ok := s.seasonCountCache[grandparentRatingKey][season]; ok {
		return cached, nil
	}

	seasonKey, err := s.findSeasonRatingKey(ctx, grandparentRatingKey, season)
	if err != nil {
		return 0, err
	}

	url, err := s.appendXToken(fmt.Sprintf("/library/metadata/%s/children", *seasonKey))
	if err != nil {
		return 0, err
	}

	raw, err := s.plexGET(ctx, *url)
	if err != nil {
		return 0, err
	}

	var res model.PlexLabels
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return 0, err
	}

	count := len(res.MediaContainer.Metadata)
	s.seasonCountCache[grandparentRatingKey][season] = count

	return count, nil
}
