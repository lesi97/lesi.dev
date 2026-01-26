package store

import (
	"context"
	"path"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

func (s *Store) UpdatePlexLabels(ctx context.Context, data model.PlexWebhookPayload) error {
	if len(data.Metadata.Location) == 0 {
		return nil
	}

	ratingKey := data.Metadata.RatingKey
	for _, location := range data.Metadata.Location {
		dir := path.Base(location.Path)
		labels, tags := s.PlexUtils.ExtractLabelsFromDir(dir)
		if len(labels) > 0 {
			err := s.PlexUtils.ApplyPlexLabels(ctx, ratingKey, "label", labels)
			if err != nil {
				return err
			}
		}

		if len(tags) > 0 {
			err := s.PlexUtils.ApplyPlexLabels(ctx, ratingKey, "genre", tags)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}
