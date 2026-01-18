package utils

type cachedWeaponPerks struct {
	CachedAtUnix int64               `json:"cached_at_unix"`
	Value        filteredPerksResult `json:"value"`
}
