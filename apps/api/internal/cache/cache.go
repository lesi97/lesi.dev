package cache

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func key(prefix string, raw string) string {
	sum := sha1.Sum([]byte(raw))
	return prefix + ":" + hex.EncodeToString(sum[:])
}

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

func (c *Cache) SetJSON(ctx context.Context, k string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, k, b, ttl).Err()
}

func (c *Cache) Key(prefix string, raw string) string {
	return key(prefix, raw)
}

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