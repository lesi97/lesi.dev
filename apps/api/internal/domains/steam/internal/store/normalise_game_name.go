package store

import "strings"

func NormaliseGameName(gameName string) string {
	trimmed := strings.TrimSpace(gameName)
	if trimmed == "" {
		return ""
	}

	return strings.ToLower(strings.Join(strings.Fields(trimmed), " "))
}
