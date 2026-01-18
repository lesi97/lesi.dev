package utils

import "strings"

func (s *Store) NormalisePayloadString(payload string) string {
	str := strings.TrimSpace(payload)

	if strings.HasPrefix(str, `"{`) || strings.HasPrefix(str, `"[\{`) {
		return str
	}

	if strings.HasPrefix(str, `\{`) {
		str = strings.TrimPrefix(str, `\`)
		str = strings.ReplaceAll(str, `\"`, `"`)
		return str
	}

	if strings.HasPrefix(str, `\[`) {
		str = strings.TrimPrefix(str, `\`)
		str = strings.ReplaceAll(str, `\"`, `"`)
		return str
	}

	return str
}
