package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func GetWeapon(ctx context.Context, database *db.DB, logger *utils.Logger, redis *redis.Client, weaponHashID string, perkHashIDs []string) (*WeaponResult, error) {
	defer logger.LogExecutionTime("MATRIX: getWeapon", time.Now(), ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan WeaponResult, 2)

	go func() {
		weapon, err := getWeaponData(ctx, database, logger, redis, weaponHashID)
		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil {
				fmt.Printf("ERROR: mtx_getWeapon: getWeaponData: %v\n", err)
			}
			ch <- WeaponResult{err: err}
			return
		}
		ch <- WeaponResult{weaponData: weapon}
	}()

	go func() {
		perks, err := getWeaponPerks(ctx, database, logger, redis, perkHashIDs)
		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil {
				fmt.Printf("ERROR: mtx_getWeapon: getWeaponPerks: %v\n", err)
			}
			ch <- WeaponResult{err: err}
			return
		}
		ch <- WeaponResult{weaponPerks: perks}
	}()

	var final WeaponResult

	for range 2 {
		res := <-ch
		if res.err != nil {
			cancel()
			return nil, res.err
		}
		if res.weaponData != nil {
			final.weaponData = res.weaponData
		}
		if res.weaponPerks != nil {
			final.weaponPerks = res.weaponPerks
		}
	}

	return &final, nil
}
