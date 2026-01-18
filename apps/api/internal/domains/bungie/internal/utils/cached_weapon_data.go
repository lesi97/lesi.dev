package utils

type cachedWeaponData struct {
	CachedAtUnix int64      `json:"cached_at_unix"`
	Value        weaponData `json:"value"`
}
