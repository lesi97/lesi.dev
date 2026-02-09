package requestmetrics

import (
	"context"
	"time"
)

func GetFetchCallsDuration(ctx context.Context) time.Duration {
	value := ctx.Value(requestMetricsContextKey{})
	if value == nil {
		return 0
	}

	metrics, ok := value.(*RequestMetrics)
	if !ok {
		return 0
	}

	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	return metrics.fetchCallsDuration
}
