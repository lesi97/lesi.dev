package store

import "time"

const steamGameNotFoundCacheValue = "__NOT_FOUND__"

var steamGameNameCacheTTL = 30 * 24 * time.Hour
var steamGameNegativeCacheTTL = 6 * time.Hour
