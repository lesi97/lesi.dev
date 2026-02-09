package requestmetrics

import (
	"context"
	"errors"
	"time"
)

func AddDatabaseCallsDuration(ctx context.Context, duration time.Duration, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return
	}
	if ctx != nil && ctx.Err() != nil {
		return
	}

	value := ctx.Value(requestMetricsContextKey{})
	if value == nil {
		return
	}

	metrics, ok := value.(*RequestMetrics)
	if !ok {
		return
	}

	metrics.lock.Lock()
	defer metrics.lock.Unlock()
	metrics.databaseCallsDuration += duration
}
