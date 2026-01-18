package utils

type cachedProfile struct {
	CachedAtUnix int64         `json:"cached_at_unix"`
	Value        BungieProfile `json:"value"`
}
