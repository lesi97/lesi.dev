package bungie_store

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
)


type specificUserWeaponData struct {
	Weapon 			string
	HashID 			string
	CrucibleTracker string
}


func (s *BungieStore) GetTerrorWeapon(ctx context.Context) (*string, error) {
	context, ok := ctx.Value("bungie").(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	weaponData, err := getTerrorWeaponData(context.WeaponName)
	if err != nil {
		return nil, err
	}

	const membershipID = "4611686018467358417"
	const preferredPlatform = "3"

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


func getTerrorWeaponData(name string) (*specificUserWeaponData, error) {
	switch name {
        case "ace":
            return  &specificUserWeaponData{Weapon: "6917529207684081719", HashID: "38912240", CrucibleTracker: "38912240"}, nil
        case "felwinter":
            return  &specificUserWeaponData{Weapon: "6917529190261952418", HashID: "1179141605", CrucibleTracker: "3244015567"}, nil
        case "matador":
            return  &specificUserWeaponData{Weapon: "6917529875871677239", HashID: "2563012876", CrucibleTracker: "3244015567"}, nil
        case "immortal":
            return  &specificUserWeaponData{Weapon: "6917529880229623656", HashID: "38912240", CrucibleTracker: "38912240"}, nil
        case "thorn":
            return &specificUserWeaponData{Weapon: "6917529935035554307", HashID: "3973202132", CrucibleTracker: "38912240"}, nil
        default:
            return nil, fmt.Errorf("weapon '%s' not found", name)
    }
}