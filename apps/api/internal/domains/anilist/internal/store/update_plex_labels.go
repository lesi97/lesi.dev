package store

import (
	"context"
	"path"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

func (s *Store) UpdatePlexLabels(ctx context.Context, data model.PlexWebhookPayload) error {
	dir := path.Base(data.Metadata.Location[0].Path)
	labels, tags := s.PlexUtils.ExtractLabelsFromDir(dir)
	ratingKey := data.Metadata.RatingKey

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
	
	return nil
}