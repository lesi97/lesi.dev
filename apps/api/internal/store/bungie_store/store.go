package bungie_store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/lesi97/lesi.dev/api/internal/database"
	"github.com/lesi97/lesi.dev/api/internal/utils"
)

// const bungie_url = "https://www.bungie.net"

type BungieStore interface {
	GetCharacterPlayTime(ctx context.Context) (*string, error)
	GetEquippedWeapon(ctx context.Context) (*string, error)
	GetTerrorWeapon(ctx context.Context) (*string, error)
}

type BungieContextInfo struct {
	Platform 	string
	Gamertag 	string
	Handler 	string
	WeaponIndex int
	WeaponName  string
}

type SupabaseBungieStore struct {
	db *database.Supabase
	logger *utils.Logger
	url string
}

func NewSupabaseBungieStore(db *database.Supabase, logger *utils.Logger) *SupabaseBungieStore {
	return &SupabaseBungieStore{
		db: db,
		logger: logger,
		url: "https://www.bungie.net",
	}
}

func (s *SupabaseBungieStore) GetCharacterPlayTime(ctx context.Context) (*string, error) {
	context, ok := ctx.Value("bungie").(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	user, err := s.getUser(ctx, context.Gamertag)
	if err != nil {
		s.logger.Printf("ERROR: getCharacterPlayTime %v\n", err)
		return nil, err
	}
	preferredPlatform := getPlatformEnum(context.Platform, user)
	characters, err := s.getBungieProfileByMembershipID(user.MembershipID, preferredPlatform, "200")
	if err != nil {
		s.logger.Printf("ERROR: getBungieProfileByMembershipID %v\n", err)
		return nil, err
	}

	var parts []string

	for _, character := range characters.Response.Characters.Data {
		class := getCharacterType(character.ClassType)
		minutes, err :=  strconv.Atoi(character.MinutesPlayedTotal)
		if err != nil {
			return nil, err
		}

		playTime := formatPlayTime(minutes)

		var formatted string
		switch {
		case playTime.Hours == 0:
			formatted = fmt.Sprintf("%s: %d minute%s", class, playTime.Minutes, plural(playTime.Minutes))
		case playTime.Minutes == 0:
			formatted = fmt.Sprintf("%s: %d hour%s", class, playTime.Hours, plural(playTime.Hours))
		default:
			formatted = fmt.Sprintf("%s: %d hour%s & %d minute%s", class, playTime.Hours, plural(playTime.Hours), playTime.Minutes, plural(playTime.Minutes))
		}

		parts = append(parts, formatted)
	}

	message := fmt.Sprintf("%s - %s", context.Gamertag, strings.Join(parts, " | "))
	return &message, nil

}

func (s *SupabaseBungieStore) GetEquippedWeapon(ctx context.Context) (*string, error) {
	context, ok := ctx.Value("bungie").(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	user, err := s.getUser(ctx, context.Gamertag)
	if err != nil {
		s.logger.Printf("ERROR in %v:\n - GetEquippedPrimary\n   - getUser: %v\n", context.Handler, err)
		return nil, err
	}

	preferredPlatform := getPlatformEnum(context.Platform, user)
	profile, err := s.getBungieProfileByMembershipID(user.MembershipID, preferredPlatform, "200,205,302,305,309")
	if err != nil {
		fmt.Printf("ERROR: GetEquippedPrimary: getBungieProfileByMembershipID: %v\n", err)
		return nil, err
	}

	mainCharId := findLastPlayedChar(s, profile.Response.Characters).CharacterID

	itemInstanceID := profile.Response.CharacterEquipment.Data[mainCharId].Items[context.WeaponIndex].ItemInstanceID
	itemHashID := strconv.Itoa(profile.Response.CharacterEquipment.Data[mainCharId].Items[context.WeaponIndex].ItemHash)

	plugHashes := profile.Response.ItemComponents.Sockets.Data[itemInstanceID].Sockets
	perkHashIDs := getPerkHashIDs(plugHashes)
	
	plugObjectives := profile.Response.ItemComponents.PlugObjectives.Data[itemInstanceID].ObjectivesPerPlug
	killCount, category := getKillCounts(plugObjectives)

	weapon, err := s.getWeapon(ctx, itemHashID, perkHashIDs)
	if err != nil {
		fmt.Printf("ERROR: GetEquippedWeapon: getWeaponData: %v\n", err)
		go func() {
			s.getNewWeapons()
		}()
		return nil, fmt.Errorf("weapon not found, please try again shortly, I'm updating my records 🤓")
	}

	go func() {
		data := dbKillCounts{
			MembershipID: user.MembershipID,
			WeaponID: itemInstanceID,
			WeaponHash: itemHashID,
		}
		switch category {
			case "PVP":
				data.PVPKills = &killCount
			case "PVE":
				data.PVEKills = &killCount
			case "Trials":
				data.TrialsKills = &killCount
		}
			
		s.insertKillCounts(&data)
	}()

	responseMessage := generateString(context.Gamertag, weapon, category, killCount)
	return &responseMessage, nil
}

func (s *SupabaseBungieStore) GetTerrorWeapon(ctx context.Context) (*string, error) {
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
		Err error
	})
	bungieChan := make(chan struct {
		Bungie *BungieProfile
		Err error
	})

	go func() {
		count, err := s.getKillCountsFromDB(ctx, membershipID, weaponData.Weapon)
		if err != nil {
			s.logger.Fatalf("ERROR: %v\n - getKillCountsFromDB: %v", context.Handler, err)
			dbChan <- struct{Count int; Err error}{Count: 0, Err: err}
			return
		}
		dbChan <- struct{Count int; Err error}{Count: count.PVPKills, Err: nil}
	}()

	go func() {
		profile, err := s.getBungieProfileByMembershipID(membershipID, preferredPlatform, "205,309")
		if err != nil {
			fmt.Printf("ERROR: %s: getBungieProfileByMembershipID: %v\n", context.Handler, err)
			bungieChan <- struct{Bungie *BungieProfile; Err error}{Bungie: nil, Err: err}
			return
		}
		bungieChan <- struct{Bungie *BungieProfile; Err error}{Bungie: profile, Err: nil}
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
