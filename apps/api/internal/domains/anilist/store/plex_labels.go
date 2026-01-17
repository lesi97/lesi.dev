package anilist_store

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (s *AnilistStore) plexLabels(ratingKey string) ([]string, error) {
	url := fmt.Sprintf(
		"%s/library/metadata/%s?X-Plex-Token=%s",
		s.plex.baseUrl,
		ratingKey,
		s.plex.xtoken,
	)

	raw, err := s.plexGet(url)
	if err != nil {
		return nil, err
	}

	var res struct {
		MediaContainer struct {
			Metadata []struct {
				Tag []struct {
					Tag     string `json:"tag"`
					TagType string `json:"tagType"`
				} `json:"Tag"`
			} `json:"Metadata"`
		} `json:"MediaContainer"`
	}

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

func hasBlockedLabel(labels []string, blocked map[string]bool) bool {
	for _, l := range labels {
		if blocked[strings.ToLower(l)] {
			return true
		}
	}
	return false
}
