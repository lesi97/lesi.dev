package utils

import "github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"

type cachedTrialsReport struct {
	CachedAtUnix int64            `json:"cached_at_unix"`
	Value        model.TrialsData `json:"value"`
}
