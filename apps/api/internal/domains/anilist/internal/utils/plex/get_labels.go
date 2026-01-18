package plex

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

func (s *Store) getLabels(ctx context.Context, ratingKey string) ([]string, error) {
	url, err := s.appendXToken(fmt.Sprintf("/library/metadata/%s", ratingKey))
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

	labels := make([]string, 0)
	if len(res.MediaContainer.Metadata) == 0 {
		return labels, nil
	}

	for _, t := range res.MediaContainer.Metadata[0].Tag {
		if t.TagType == "label" {
			labels = append(labels, t.Tag)
		}
	}

	return labels, nil
}
