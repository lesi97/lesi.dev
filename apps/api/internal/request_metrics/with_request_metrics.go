package requestmetrics

import "context"

func WithRequestMetrics(ctx context.Context) context.Context {
	metrics := &RequestMetrics{}
	return context.WithValue(ctx, requestMetricsContextKey{}, metrics)
}
