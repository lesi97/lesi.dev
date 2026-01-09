package anilist_store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type plexChildrenResponse struct {
	MediaContainer struct {
		Metadata []struct {
			RatingKey string `json:"ratingKey"`
			Index     int    `json:"index"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

func (s *AnilistStore) plexGet(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("plex http %d %s", resp.StatusCode, string(raw))
	}

	return raw, nil
}

func (s *AnilistStore) findSeasonRatingKey(grandparentRatingKey string, season int) (string, error) {
	if strings.TrimSpace(s.plex.baseUrl) == "" || strings.TrimSpace(s.plex.xtoken) == "" {
		return "", fmt.Errorf("plex base url or token missing")
	}

	url := fmt.Sprintf("%s/library/metadata/%s/children?X-Plex-Token=%s", s.plex.baseUrl, grandparentRatingKey, s.plex.xtoken)
	raw, err := s.plexGet(url)
	if err != nil {
		return "", err
	}

	var res plexChildrenResponse
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return "", err
	}

	for _, m := range res.MediaContainer.Metadata {
		if m.Index == season {
			return m.RatingKey, nil
		}
	}

	return "", fmt.Errorf("season %d not found for show %s", season, grandparentRatingKey)
}

func (s *AnilistStore) seasonEpisodeCount(grandparentRatingKey string, season int) (int, error) {
	if _, ok := s.seasonCountCache[grandparentRatingKey]; !ok {
		s.seasonCountCache[grandparentRatingKey] = make(map[int]int)
	}
	if cached, ok := s.seasonCountCache[grandparentRatingKey][season]; ok {
		return cached, nil
	}

	seasonKey, err := s.findSeasonRatingKey(grandparentRatingKey, season)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s/library/metadata/%s/children?X-Plex-Token=%s", s.plex.baseUrl, seasonKey, s.plex.xtoken)
	raw, err := s.plexGet(url)
	if err != nil {
		return 0, err
	}

	var res plexChildrenResponse
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return 0, err
	}

	count := len(res.MediaContainer.Metadata)
	s.seasonCountCache[grandparentRatingKey][season] = count

	return count, nil
}

func (s *AnilistStore) absoluteEpisodeFromPlex(grandparentRatingKey string, season int, episode int) (int, error) {
	if season <= 1 {
		return episode, nil
	}
	if season == 0 {
		return 0, fmt.Errorf("season 0 specials should not be converted to absolute episode")
	}

	total := 0
	for sn := 1; sn < season; sn++ {
		c, err := s.seasonEpisodeCount(grandparentRatingKey, sn)
		if err != nil {
			return 0, err
		}
		total += c
	}

	return total + episode, nil
}
