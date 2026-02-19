package store

import (
	"strings"
	"unicode"
)

func NormaliseGameName(gameName string) string {
	if strings.TrimSpace(gameName) == "" {
		return ""
	}

	mapped := strings.Map(func(char rune) rune {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			return unicode.ToLower(char)
		}
		if unicode.IsSpace(char) {
			return ' '
		}
		return ' '
	}, gameName)

	return strings.Join(strings.Fields(mapped), " ")
}
