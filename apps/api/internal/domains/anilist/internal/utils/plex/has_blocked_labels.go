package plex

import "strings"

func hasBlockedLabel(labels []string, blockedLabels map[string]bool) bool {
	for _, l := range labels {
		if blockedLabels[strings.ToLower(l)] {
			return true
		}
	}
	return false
}
