package store

import (
	"context"
	"fmt"
	"strconv"

	bungie_utils "github.com/lesi97/lesi.dev/internal/domains/bungie/internal/utils"
)

func (s *Store) GetEquippedWeapon(ctx context.Context) (*string, error) {
	request, ok := ctx.Value(bungieContextKey).(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}

	user, err := bungie_utils.GetUser(ctx, s.DB, s.Logger, s.Redis, s.BaseURL, s.ClientID, request.Gamertag, request.Platform)
	if err != nil {
		s.Logger.Printf("ERROR in %v:\n - GetEquippedPrimary\n   - getUser: %v\n", request.Handler, err)
		return nil, err
	}

	preferredPlatform := bungie_utils.GetPlatformEnum(request.Platform, user.MembershipType)
	profile, err := bungie_utils.GetBungieProfileByMembershipID(
		ctx,
		s.Redis,
		s.Logger,
		s.ClientID,
		s.BaseURL,
		user.MembershipID,
		preferredPlatform,
		"200,205,302,305,309",
	)
	if err != nil {
		fmt.Printf("ERROR: GetEquippedPrimary: getBungieProfileByMembershipID: %v\n", err)
		return nil, err
	}

	mainChar := bungie_utils.FindLastPlayedCharacter(s.DB, s.Logger, profile.Response.Characters)
	if mainChar == nil {
		return nil, fmt.Errorf("no characters found")
	}
	mainCharId := mainChar.CharacterID

	itemInstanceID := profile.Response.CharacterEquipment.Data[mainCharId].Items[request.WeaponIndex].ItemInstanceID
	itemHashID := strconv.Itoa(profile.Response.CharacterEquipment.Data[mainCharId].Items[request.WeaponIndex].ItemHash)

	plugHashes := profile.Response.ItemComponents.Sockets.Data[itemInstanceID].Sockets
	perkHashIDs := bungie_utils.GetPerkHashIDs(plugHashes)

	plugObjectives := profile.Response.ItemComponents.PlugObjectives.Data[itemInstanceID].ObjectivesPerPlug
	killCount, category := bungie_utils.GetKillCounts(plugObjectives)

	weapon, err := bungie_utils.GetWeapon(ctx, s.DB, s.Logger, s.Redis, itemHashID, perkHashIDs)
	if err != nil {
		fmt.Printf("ERROR: GetEquippedWeapon: getWeaponData: %v\n", err)
		go func() {
			bungie_utils.GetNewWeapons(s.DB, s.Logger, s.BaseURL, s.ClientID)
		}()
		return nil, fmt.Errorf("weapon not found, please try again shortly, I'm updating my records ??")
	}

	go func() {
		switch category {
		case "PVP":
			bungie_utils.InsertKillCounts(s.DB, s.Logger, user.MembershipID, itemInstanceID, itemHashID, &killCount, nil, nil)
		case "PVE":
			bungie_utils.InsertKillCounts(s.DB, s.Logger, user.MembershipID, itemInstanceID, itemHashID, nil, &killCount, nil)
		case "Trials":
			bungie_utils.InsertKillCounts(s.DB, s.Logger, user.MembershipID, itemInstanceID, itemHashID, nil, nil, &killCount)
		default:
			bungie_utils.InsertKillCounts(s.DB, s.Logger, user.MembershipID, itemInstanceID, itemHashID, nil, nil, nil)
		}
	}()

	responseMessage := bungie_utils.GenerateString(request.Gamertag, weapon, category, killCount)
	return &responseMessage, nil
}
