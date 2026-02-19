package store

import "time"

const steamGameNotFoundCacheValue = "__NOT_FOUND__"

var steamGameNameCacheTTL = 90 * 24 * time.Hour
var steamGameNegativeCacheTTL = 30 * time.Minute
