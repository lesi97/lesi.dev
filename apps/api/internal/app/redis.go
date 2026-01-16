package app

import (
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func RedisLoadConfig() Config {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	return Config{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
}

func RedisNew(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    2 * time.Second,
		MinRetryBackoff: 50 * time.Millisecond,
		MaxRetryBackoff: 250 * time.Millisecond,
		MaxRetries:      2,
	})
}
