package requestmetrics

import (
	"sync"
	"time"
)

type requestMetricsContextKey struct{}

type RequestMetrics struct {
	lock                  sync.Mutex
	fetchCallsDuration    time.Duration
	databaseCallsDuration time.Duration
}
