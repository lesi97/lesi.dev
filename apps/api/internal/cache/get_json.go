package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func (c *Cache) GetJSON(ctx context.Context, k string, dst any) (bool, error) {
	s, err := c.rdb.Get(ctx, k).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(s), dst); err != nil {
		_ = c.rdb.Del(ctx, k).Err()
		return false, nil
	}

	return true, nil
}