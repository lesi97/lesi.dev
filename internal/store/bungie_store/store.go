package bungie_store

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lesi97/api.lesi.dev/internal/database"
)

const bungie_url = "https://www.bungie.net"

type BungieStore interface {
	GetCharacterPlayTime(ctx context.Context) (*string, error)
	GetEquippedWeapon(ctx context.Context) (*string, error)
}

type BungieContextInfo struct {
	Platform 	string
	Gamertag 	string
	Handler 	string
	WeaponIndex int
}

type SupabaseBungieStoreStore struct {
	db *database.Supabase
	logger *log.Logger
}

func NewSupabaseBungieStore(db *database.Supabase, logger *log.Logger) *SupabaseBungieStoreStore {
	return &SupabaseBungieStoreStore{
		db: db,
		logger: logger,
	}
}

func (s *SupabaseBungieStoreStore) GetCharacterPlayTime(ctx context.Context) (*string, error) {
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
	characters, err := getBungieProfileByMembershipID(user.MembershipID, preferredPlatform, "200")
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

func (s *SupabaseBungieStoreStore) GetEquippedWeapon(ctx context.Context) (*string, error) {
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
	profile, err := getBungieProfileByMembershipID(user.MembershipID, preferredPlatform, "200,205,302,305,309")
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
		return nil, err
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