package cache

import (
	"context"
	"time"
)

func (c *Cache) WithLock(ctx context.Context, lockKey string, ttl time.Duration, fn func() error) error {
	ok, err := c.rdb.SetNX(ctx, lockKey, "1", ttl).Result()
	if err != nil {
		return err
	}

	if ok {
		defer c.rdb.Del(ctx, lockKey).Err()
		return fn()
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
		exists, err := c.rdb.Exists(ctx, lockKey).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			return nil
		}
	}

	return nil
}