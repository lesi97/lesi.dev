package bungie_store

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

// const bungie_url = "https://www.bungie.net"

type BungieStoreInterface interface {
	GetCharacterPlayTime(ctx context.Context) (*string, error)
	GetEquippedWeapon(ctx context.Context) (*string, error)
	GetTerrorWeapon(ctx context.Context) (*string, error)
	GetYungerWeapon(ctx context.Context) (*string, error)
}

type BungieContextInfo struct {
	Platform 	string
	Gamertag 	string
	Handler 	string
	WeaponIndex int
	WeaponName  string
}

type BungieStore struct {
	store.StoreBase
	redis 		*redis.Client
	url 		string
	clientId 	string
}

func NewStore(db *database.DB, logger *utils.Logger, redis *redis.Client) (*BungieStore, error) {
	bungieApiKey := os.Getenv("BUNGIE_CLIENT_ID")
	if bungieApiKey == "" {
		message := "FATAL: ERROR GETTING BUNGIE_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: "BUNGIE STORE FATAL",
			Title: "BUNGIE STORE FATAL",
		})
		return nil, fmt.Errorf("%s", message)
	}

	return &BungieStore{
		StoreBase: store.NewStoreBase(db, logger),
		redis: redis,
		url: "https://www.bungie.net",
		clientId: bungieApiKey,
	}, nil
}

func (s *BungieStore) GetCharacterPlayTime(ctx context.Context) (*string, error) {
	context, ok := ctx.Value("bungie").(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	user, err := s.getUser(ctx, context.Gamertag, context.Platform)
	if err != nil {
		s.Logger.Printf("ERROR: getCharacterPlayTime %v\n", err)
		return nil, err
	}
	preferredPlatform := getPlatformEnum(context.Platform, user)
	characters, err := s.getBungieProfileByMembershipID(user.MembershipID, preferredPlatform, "200")
	if err != nil {
		s.Logger.Printf("ERROR: getBungieProfileByMembershipID %v\n", err)
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

func (s *BungieStore) GetEquippedWeapon(ctx context.Context) (*string, error) {
	context, ok := ctx.Value("bungie").(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}
	user, err := s.getUser(ctx, context.Gamertag, context.Platform)
	if err != nil {
		s.Logger.Printf("ERROR in %v:\n - GetEquippedPrimary\n   - getUser: %v\n", context.Handler, err)
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

