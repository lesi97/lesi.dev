package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func getWeaponData(ctx context.Context, database *db.DB, logger *utils.Logger, redis *redis.Client, hashID string) (*weaponData, error) {
	startedAt := time.Now()
	shouldLog := true
	defer func() {
		if shouldLog {
			logger.LogExecutionTime("DATABASE CALL: getWeaponData", startedAt, ctx)
		}
	}()
	freshFor := 7 * 24 * time.Hour
	staleFor := 30 * 24 * time.Hour
	now := time.Now()
	cacheKey := fmt.Sprintf("bungie:weapon:data:%s", hashID)

	cached, err := redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrap cachedWeaponData
		if err := json.Unmarshal([]byte(cached), &wrap); err == nil {
			age := now.Sub(time.Unix(wrap.CachedAtUnix, 0))
			if age <= freshFor {
				logger.PrintColour(true, "brightBlack", "CACHE HIT fresh getWeaponData %s", cacheKey)
				shouldLog = false
				v := wrap.Value
				return &v, nil
			}
			if age <= staleFor {
				logger.PrintColour(true, "brightBlack", "CACHE HIT stale getWeaponData %s", cacheKey)
				shouldLog = false
				lockKey := cacheKey + ":lock"
				ok, _ := redis.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
				if ok {
					go func() {
						defer redis.Del(context.Background(), lockKey).Err()
						data := &weaponData{HashID: hashID}
						query := `
		SELECT 
			display_name, 
			tier_type_name
		FROM destiny_weapons 
		WHERE id = $1
	`
						err := database.QueryRow(context.Background(), query, hashID).Scan(
							&data.DisplayName,
							&data.TierTypeName,
						)
						if err != nil || err == sql.ErrNoRows {
							return
						}
						wrap := cachedWeaponData{
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
			display_name, 
			tier_type_name
		FROM destiny_weapons 
		WHERE id = $1
	`

	data := &weaponData{
		HashID: hashID,
	}

	err = database.QueryRow(ctx, query, hashID).Scan(
		&data.DisplayName,
		&data.TierTypeName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	wrap := cachedWeaponData{
		CachedAtUnix: time.Now().Unix(),
		Value:        *data,
	}
	if b, err := json.Marshal(wrap); err == nil {
		_ = redis.Set(ctx, cacheKey, b, staleFor).Err()
	}

	return data, nil
}
