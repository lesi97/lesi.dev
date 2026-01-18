package anilist

import "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"

func PickBestMatch(list []model.AnilistMedia, showName string, plexYear int) (*model.AnilistMedia, bool) {
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
