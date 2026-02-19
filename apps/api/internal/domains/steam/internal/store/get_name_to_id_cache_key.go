package store

import (
	"fmt"
	"net/url"
)

func GetNameToIDCacheKey(normalisedGameName string) string {
	return fmt.Sprintf("steam:game:name_to_id:%s", url.QueryEscape(normalisedGameName))
}
