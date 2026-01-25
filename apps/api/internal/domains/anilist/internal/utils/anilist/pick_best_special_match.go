package anilist

import (
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
)

func PickBestSpecialMatch(list []model.AnilistMedia, showName string, episodeTitle string, plexYear int) (*model.AnilistMedia, bool) {
	if len(list) == 0 {
		return nil, false
	}

	bestIdx := 0
	bestScore := 0
	needleShow := strings.ToLower(strings.TrimSpace(showName))
	needleEpisode := strings.ToLower(strings.TrimSpace(episodeTitle))

	for i := range list {
		m := list[i]
		score := 0

		titles := []string{m.Title.Romaji, m.Title.English, m.Title.Native}
		for _, title := range titles {
			if strings.TrimSpace(title) == "" {
				continue
			}

			trimmed := strings.TrimSpace(title)
			lower := strings.ToLower(trimmed)

			if needleShow != "" && lower == needleShow {
				score += 400
			}
			if needleEpisode != "" && lower == needleEpisode {
				score += 800
			}
			if needleShow != "" && strings.Contains(lower, needleShow) {
				score += 200
			}
			if needleEpisode != "" && strings.Contains(lower, needleEpisode) {
				score += 600
			}
			if needleShow != "" && needleEpisode != "" && strings.Contains(lower, needleShow) && strings.Contains(lower, needleEpisode) {
				score += 1000
			}
		}

		if plexYear > 0 && m.SeasonYear != nil {
			diff := *m.SeasonYear - plexYear
			if diff < 0 {
				diff = -diff
			}
			if diff == 0 {
				score += 200
			} else if diff <= 1 {
				score += 100
			} else if diff <= 2 {
				score += 50
			} else {
				score -= diff
			}
		}

		if m.IsAdult {
			score -= 10000
		}

		if i == 0 || score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	return &list[bestIdx], true
}
