package plex

import (
	"context"
	"fmt"
)

func (s *Store) AbsoluteEpisodeFromPlex(ctx context.Context, grandparentRatingKey string, season int, episode int) (int, error) {
	if season <= 1 {
		return episode, nil
	}
	if season == 0 {
		return 0, fmt.Errorf("season 0 specials should not be converted to absolute episode")
	}

	total := 0
	for sn := 1; sn < season; sn++ {
		c, err := s.getSeasonEpisodeCount(ctx, grandparentRatingKey, sn)
		if err != nil {
			return 0, err
		}
		total += c
	}

	return total + episode, nil
}
