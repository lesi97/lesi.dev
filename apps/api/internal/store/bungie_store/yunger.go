package bungie_store

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
)

func (s *BungieStore) GetYungerWeapon(ctx context.Context) (*string, error) {
	context, ok := ctx.Value("bungie").(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	weaponData, err := getYungerWeaponData(context.WeaponName)
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
		Bungie *BungieProfile
		Err    error
	})

	go func() {
		count, err := s.getKillCountsFromDB(ctx, membershipID, weaponData.Weapon)
		if err != nil {
			s.Logger.Fatalf("ERROR: %v\n - getKillCountsFromDB: %v", context.Handler, err)
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
		profile, err := s.getBungieProfileByMembershipID(membershipID, preferredPlatform, "205,309")
		if err != nil {
			fmt.Printf("ERROR: %s: getBungieProfileByMembershipID: %v\n", context.Handler, err)
			bungieChan <- struct {
				Bungie *BungieProfile
				Err    error
			}{Bungie: nil, Err: err}
			return
		}
		bungieChan <- struct {
			Bungie *BungieProfile
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
				count, category := getKillCounts(objectives)
				killCount = count
				valid = true

				go func() {
					data := dbKillCounts{
						MembershipID: membershipID,
						WeaponID:     weaponData.Weapon,
						WeaponHash:   weaponData.HashID,
					}
					switch category {
					case "PVP":
						data.PVPKills = &count
					case "PVE":
						data.PVEKills = &count
					case "Trials":
						data.TrialsKills = &count
					}
					s.insertKillCounts(&data)
				}()
			}
		}
	}

	localeKillCount := humanize.Comma(int64(killCount))
	return &localeKillCount, nil

}



func getYungerWeaponData(name string) (*specificUserWeaponData, error) {
	switch name {
        case "cloudstrike":
            return  &specificUserWeaponData{Weapon: "6917529926351690059", HashID: "", CrucibleTracker: "3244015567"}, nil
        case "beloved":
            return  &specificUserWeaponData{Weapon: "6917529798698579887", HashID: "", CrucibleTracker: "3244015567"}, nil
        default:
            return nil, fmt.Errorf("weapon '%s' not found", name)
    }
}