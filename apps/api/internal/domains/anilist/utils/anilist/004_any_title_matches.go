package anilist

import "github.com/lesi97/lesi.dev/internal/domains/anilist/model"

func anyTitleMatches(m model.AnilistMedia, showName string) bool {
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
