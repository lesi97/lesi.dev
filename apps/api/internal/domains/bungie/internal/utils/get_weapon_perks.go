package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/cache"
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func getWeaponPerks(ctx context.Context, database *db.DB, logger *utils.Logger, redis *redis.Client, perkHashIDs []string) (*filteredPerksResult, error) {
	startedAt := time.Now()
	shouldLog := true
	defer func() {
		if shouldLog {
			logger.LogExecutionTime("DATABASE CALL: getWeaponPerks", startedAt, ctx)
		}
	}()
	if len(perkHashIDs) == 0 {
		return nil, fmt.Errorf("perk list not provided")
	}

	freshFor := 7 * 24 * time.Hour
	staleFor := 30 * 24 * time.Hour
	now := time.Now()
	cacheKey := fmt.Sprintf("bungie:weapon:perks:%s", strings.Join(perkHashIDs, ","))

	cached, err := redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrap cachedWeaponPerks
		if err := json.Unmarshal([]byte(cached), &wrap); err == nil {
			age := now.Sub(time.Unix(wrap.CachedAtUnix, 0))
			if age <= freshFor {
				logger.PrintCache("CACHE HIT fresh getWeaponPerks %s", cacheKey)
				shouldLog = false
				v := wrap.Value
				return &v, nil
			}
			if age <= staleFor {
				logger.PrintCache("CACHE HIT stale getWeaponPerks %s", cacheKey)
				shouldLog = false
				lockKey := cacheKey + ":lock"
				ok, _ := redis.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
				if ok {
					go func() {
						defer redis.Del(context.Background(), lockKey).Err()
						placeholders := make([]string, len(perkHashIDs))
						args := make([]interface{}, len(perkHashIDs))
						for i, id := range perkHashIDs {
							placeholders[i] = fmt.Sprintf("$%d", i+1)
							args[i] = id
						}
						query := fmt.Sprintf(`
		select 
			name, 
			item_type 
		from destiny_weapon_perks 
		where hash_id in (%s)
	`, strings.Join(placeholders, ", "))

						rows, err := database.Query(context.Background(), query, args...)
						if err != nil {
							return
						}
						defer rows.Close()

						perks := make([]databasePerk, 0, 16)
						for rows.Next() {
							var perk databasePerk
							if err := rows.Scan(&perk.Name, &perk.ItemType); err != nil {
								return
							}
							perks = append(perks, perk)
						}
						if err := rows.Err(); err != nil {
							return
						}

						filtered := filterWeaponPerks(perks)
						wrap := cachedWeaponPerks{
							CachedAtUnix: time.Now().Unix(),
							Value:        filtered,
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
	} else if !cache.IsRedisNil(err) {
		return nil, err
	}

	placeholders := make([]string, len(perkHashIDs))
	args := make([]interface{}, len(perkHashIDs))
	for i, id := range perkHashIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`
		select 
			name, 
			item_type 
		from destiny_weapon_perks 
		where hash_id in (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := database.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perks := make([]databasePerk, 0, 16)

	for rows.Next() {
		var perk databasePerk
		if err := rows.Scan(&perk.Name, &perk.ItemType); err != nil {
			return nil, err
		}
		perks = append(perks, perk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	filtered := filterWeaponPerks(perks)
	wrap := cachedWeaponPerks{
		CachedAtUnix: time.Now().Unix(),
		Value:        filtered,
	}
	if b, err := json.Marshal(wrap); err == nil {
		_ = redis.Set(ctx, cacheKey, b, staleFor).Err()
	}
	return &filtered, nil
}
