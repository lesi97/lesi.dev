package plex

import (
	"regexp"
	"sort"
	"strings"
)

func (s *Store) ExtractLabelsFromDir(dir string) ([]string, []string) {
	re := regexp.MustCompile(`(?i)\b(labels|tags)\[([^\]]*)\]`)

	matches := re.FindAllStringSubmatch(dir, -1)

	labelSet := map[string]struct{}{}
	tagSet := map[string]struct{}{}

	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		kind := strings.ToLower(strings.TrimSpace(m[1]))
		rawList := m[2]

		for _, raw := range strings.Split(rawList, ",") {
			v := strings.TrimSpace(raw)
			if v == "" {
				continue
			}
			v = normaliseSpaces(v)

			if kind == "labels" {
				labelSet[v] = struct{}{}
			} else {
				tagSet[v] = struct{}{}
			}
		}
	}

	labels := setToSortedSlice(labelSet)
	tags := setToSortedSlice(tagSet)

	return labels, tags
}

func normaliseSpaces(s string) string {
	parts := strings.Fields(s)
	return strings.Join(parts, " ")
}

func setToSortedSlice(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}