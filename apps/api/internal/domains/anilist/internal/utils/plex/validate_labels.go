package plex

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

func (s *Store) ValidateLabels(ctx context.Context, d *model.PlexWebhookPayload) (bool, error) {

	var blockedLabels = map[string]bool{
		"no_anilist": true,
		"no_sync":    true,
		"private":    true,
	}

	episodeLabels, err := s.getLabels(ctx, d.Metadata.RatingKey)
	if err != nil {
		return false, err
	}
	if hasBlockedLabel(episodeLabels, blockedLabels) {
		return false, nil
	}

	seasonLabels, err := s.getLabels(ctx, d.Metadata.ParentRatingKey)
	if err != nil {
		return false, err
	}
	if hasBlockedLabel(seasonLabels, blockedLabels) {
		return false, nil
	}

	showLabels, err := s.getLabels(ctx, d.Metadata.GrandparentRatingKey)
	if err != nil {
		return false, err
	}
	if hasBlockedLabel(showLabels, blockedLabels) {
		return false, nil
	}

	return true, nil
}
