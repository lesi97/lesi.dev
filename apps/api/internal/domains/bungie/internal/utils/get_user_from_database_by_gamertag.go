package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type bungieDBData struct {
	BungieID          string `json:"bungie_id"`
	MembershipID      string `json:"membership_id"`
	PreferredPlatform int64  `json:"preferred_platform"`
	FriendlyName      string `json:"friendly_name"`
}

func getUserFromDatabaseByGamertag(ctx context.Context, database *db.DB, logger *utils.Logger, redis *redis.Client, bungieID string) (*bungieDBData, error) {
	startedAt := time.Now()
	shouldLog := true
	defer func() {
		if ctx.Err() == nil && shouldLog {
			logger.LogExecutionTime("DATABASE CALL: getUserFromDatabaseByGamertag", startedAt, ctx)
		}
	}()
	freshFor := 7 * 24 * time.Hour
	staleFor := 14 * 24 * time.Hour
	now := time.Now()
	cacheKey := fmt.Sprintf("bungie:user:db:%s", strings.ToLower(bungieID))

	cached, err := redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrap cachedBungieUser
		if err := json.Unmarshal([]byte(cached), &wrap); err == nil {
			age := now.Sub(time.Unix(wrap.CachedAtUnix, 0))
			if age <= freshFor {
				logger.PrintCache("CACHE HIT fresh getUserFromDatabaseByGamertag %s", cacheKey)
				shouldLog = false
				v := wrap.Value
				return &v, nil
			}
			if age <= staleFor {
				logger.PrintCache("CACHE HIT stale getUserFromDatabaseByGamertag %s", cacheKey)
				shouldLog = false
				lockKey := cacheKey + ":lock"
				ok, _ := redis.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
				if ok {
					go func() {
						defer redis.Del(context.Background(), lockKey).Err()
						query := `
		SELECT 
			membership_id, 
			preferred_platform, 
			friendly_name 
		FROM destiny_users 
		WHERE lower(bungie_id) = lower($1)
	`
						data := &bungieDBData{
							BungieID: bungieID,
						}
						err := database.QueryRow(context.Background(), query, bungieID).Scan(
							&data.MembershipID,
							&data.PreferredPlatform,
							&data.FriendlyName,
						)
						if err != nil || err == sql.ErrNoRows {
							return
						}
						wrap := cachedBungieUser{
							CachedAtUnix: time.Now().Unix(),
							Value:        *data,
						}
						if b, err := json.Marshal(wrap); err == nil {
							_ = redis.Set(context.Background(), cacheKey, b, staleFor).Err()
						}
					}()
				}
				v := wrap.Value
				return &v, nil
			}
		} else {
			_ = redis.Del(ctx, cacheKey).Err()
		}
	} else if !IsRedisNil(err) {
		return nil, err
	}

	query := `
		SELECT 
			membership_id, 
			preferred_platform, 
			friendly_name 
		FROM destiny_users 
		WHERE lower(bungie_id) = lower($1)
	`
	data := &bungieDBData{
		BungieID: bungieID,
	}
	err = database.QueryRow(ctx, query, bungieID).Scan(
		&data.MembershipID,
		&data.PreferredPlatform,
		&data.FriendlyName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	wrap := cachedBungieUser{
		CachedAtUnix: time.Now().Unix(),
		Value:        *data,
	}
	if b, err := json.Marshal(wrap); err == nil {
		_ = redis.Set(ctx, cacheKey, b, staleFor).Err()
	}

	return data, nil
}
