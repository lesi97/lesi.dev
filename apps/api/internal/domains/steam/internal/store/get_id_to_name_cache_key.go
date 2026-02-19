package store

import "fmt"

func GetIDToNameCacheKey(gameID string) string {
	return fmt.Sprintf("steam:game:id_to_name:%s", gameID)
}
