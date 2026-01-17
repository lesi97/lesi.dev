package cache

import (
	"context"
	"encoding/json"
	"time"
)

func (c *Cache) SetJSON(ctx context.Context, k string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, k, b, ttl).Err()
}