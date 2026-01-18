package anilist

import "strings"

func titleEqualsIgnoreCase(a string, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
