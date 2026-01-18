package plex

import (
	"regexp"
	"strings"
)

func normalisePlexTitle(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	punctuationRegex := regexp.MustCompile(`[:!?.'",()\[\]{}]`)
	value = punctuationRegex.ReplaceAllString(value, "")

	spaceRegex := regexp.MustCompile(`\s+`)
	value = spaceRegex.ReplaceAllString(value, " ")

	return value
}
