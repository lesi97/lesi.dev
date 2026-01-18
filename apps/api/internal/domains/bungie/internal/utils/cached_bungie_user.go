package utils

type cachedBungieUser struct {
	CachedAtUnix int64        `json:"cached_at_unix"`
	Value        bungieDBData `json:"value"`
}
