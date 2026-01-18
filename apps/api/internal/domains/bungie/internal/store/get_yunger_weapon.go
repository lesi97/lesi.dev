package store

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	bungie_utils "github.com/lesi97/lesi.dev/internal/domains/bungie/internal/utils"
)

func (s *Store) GetYungerWeapon(ctx context.Context) (*string, error) {
	request, ok := ctx.Value(bungieContextKey).(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}

	weaponData, err := bungie_utils.GetYungerWeaponData(request.WeaponName)
	if err != nil {
		return nil, err
	}

	const membershipID = "4611686018448992398"
	const preferredPlatform = "1"

	dbChan := make(chan struct {
		Count int
		Err   error
	})
	bungieChan := make(chan struct {
		Bungie *bungie_utils.BungieProfile
		Err    error
	})

	go func() {
		count, err := bungie_utils.GetKillCountsFromDB(ctx, s.DB, s.Logger, membershipID, weaponData.Weapon)
		if err != nil {
			s.Logger.Fatalf("ERROR: %v\n - getKillCountsFromDB: %v", request.Handler, err)
			dbChan <- struct {
				Count int
				Err   error
			}{Count: 0, Err: err}
			return
		}
		dbChan <- struct {
			Count int
			Err   error
		}{Count: count.PVPKills, Err: nil}
	}()

	go func() {
		profile, err := bungie_utils.GetBungieProfileByMembershipID(
			s.Redis,
			s.Logger,
			s.ClientID,
			s.BaseURL,
			membershipID,
			preferredPlatform,
			"205,309",
		)
		if err != nil {
			fmt.Printf("ERROR: %s: getBungieProfileByMembershipID: %v\n", request.Handler, err)
			bungieChan <- struct {
				Bungie *bungie_utils.BungieProfile
				Err    error
			}{Bungie: nil, Err: err}
			return
		}
		bungieChan <- struct {
			Bungie *bungie_utils.BungieProfile
			Err    error
		}{Bungie: profile, Err: nil}
	}()

	var killCount int
	var valid bool

	for range 2 {
		select {
		case dbResult := <-dbChan:
			if dbResult.Err == nil && !valid {
				killCount = dbResult.Count
			}
		case bungieResult := <-bungieChan:
			if bungieResult.Err == nil {
				profile := bungieResult.Bungie
				objectives := profile.Response.ItemComponents.PlugObjectives.Data[weaponData.Weapon].ObjectivesPerPlug
				if objectives == nil {
					continue
				}
				count, category := bungie_utils.GetKillCounts(objectives)
				killCount = count
				valid = true

				go func() {
					switch category {
					case "PVP":
						bungie_utils.InsertKillCounts(s.DB, s.Logger, membershipID, weaponData.Weapon, weaponData.HashID, &count, nil, nil)
					case "PVE":
						bungie_utils.InsertKillCounts(s.DB, s.Logger, membershipID, weaponData.Weapon, weaponData.HashID, nil, &count, nil)
					case "Trials":
						bungie_utils.InsertKillCounts(s.DB, s.Logger, membershipID, weaponData.Weapon, weaponData.HashID, nil, nil, &count)
					default:
						bungie_utils.InsertKillCounts(s.DB, s.Logger, membershipID, weaponData.Weapon, weaponData.HashID, nil, nil, nil)
					}
				}()
			}
		}
	}

	localeKillCount := humanize.Comma(int64(killCount))
	return &localeKillCount, nil
}
