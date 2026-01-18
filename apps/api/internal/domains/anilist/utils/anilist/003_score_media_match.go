package anilist

import (
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
)

func scoreMediaMatch(m model.AnilistMedia, showName string, plexYear int) int {
	score := 0

	if anyTitleMatches(m, showName) {
		score += 1000
	}

	if plexYear > 0 && m.SeasonYear != nil {
		diff := *m.SeasonYear - plexYear
		if diff < 0 {
			diff = -diff
		}
		if diff == 0 {
			score += 500
		} else if diff <= 1 {
			score += 200
		} else if diff <= 2 {
			score += 100
		} else {
			score -= diff
		}
	}

	if strings.EqualFold(m.Format, "TV") {
		score += 50
	}
	if strings.EqualFold(m.Format, "TV_SHORT") {
		score += 10
	}
	if strings.EqualFold(m.Format, "MOVIE") {
		score -= 50
	}
	if strings.EqualFold(m.Format, "OVA") {
		score -= 20
	}

	if m.IsAdult {
		score -= 10000
	}

	return score
}
