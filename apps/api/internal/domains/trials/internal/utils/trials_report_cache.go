package utils

import (
	"time"

	"sync"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
)


type TrialsReportCache struct {
	Data     *model.TrialsData
	CachedAt time.Time
}

var trialsReportCacheMu sync.Mutex
var trialsReportCache *TrialsReportCache