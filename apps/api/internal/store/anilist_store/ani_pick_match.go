package anilist_store

import (
	"strings"
)

func titleEqualsIgnoreCase(a string, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func anyTitleMatches(m mediaData, showName string) bool {
	if titleEqualsIgnoreCase(m.Title.Romaji, showName) {
		return true
	}
	if titleEqualsIgnoreCase(m.Title.English, showName) {
		return true
	}
	if titleEqualsIgnoreCase(m.Title.Native, showName) {
		return true
	}
	return false
}

func scoreMediaMatch(m mediaData, showName string, plexYear int) int {
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

func pickBestAniListMatch(list []mediaData, showName string, plexYear int) (*mediaData, bool) {
	if len(list) == 0 {
		return nil, false
	}

	bestIdx := 0
	bestScore := scoreMediaMatch(list[0], showName, plexYear)

	for i := 1; i < len(list); i++ {
		s := scoreMediaMatch(list[i], showName, plexYear)
		if s > bestScore {
			bestScore = s
			bestIdx = i
		}
	}

	return &list[bestIdx], true
}
