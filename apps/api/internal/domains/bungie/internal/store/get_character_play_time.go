package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	bungie_utils "github.com/lesi97/lesi.dev/internal/domains/bungie/internal/utils"
)

func (s *Store) GetCharacterPlayTime(ctx context.Context) (*string, error) {
	request, ok := ctx.Value(bungieContextKey).(BungieContextInfo)
	if !ok {
		return nil, fmt.Errorf("invalid context")
	}

	user, err := bungie_utils.GetUser(ctx, s.DB, s.Logger, s.Redis, s.HTTPClient, s.BaseURL, s.ClientID, request.Gamertag, request.Platform)
	if err != nil {
		s.Logger.Printf("ERROR: getCharacterPlayTime %v\n", err)
		return nil, err
	}

	preferredPlatform := bungie_utils.GetPlatformEnum(request.Platform, user.MembershipType)
	characters, err := bungie_utils.GetBungieProfileByMembershipID(ctx, s.Redis, s.Logger, s.HTTPClient, s.ClientID, s.BaseURL, user.MembershipID, preferredPlatform, "200")
	if err != nil {
		s.Logger.Printf("ERROR: getBungieProfileByMembershipID %v\n", err)
		return nil, err
	}

	var parts []string

	for _, character := range characters.Response.Characters.Data {
		class := bungie_utils.GetCharacterType(character.ClassType)
		minutes, err := strconv.Atoi(character.MinutesPlayedTotal)
		if err != nil {
			return nil, err
		}

		playTime := bungie_utils.FormatPlayTime(minutes)

		var formatted string
		switch {
		case playTime.Hours == 0:
			formatted = fmt.Sprintf("%s: %d minute%s", class, playTime.Minutes, bungie_utils.Plural(playTime.Minutes))
		case playTime.Minutes == 0:
			formatted = fmt.Sprintf("%s: %d hour%s", class, playTime.Hours, bungie_utils.Plural(playTime.Hours))
		default:
			formatted = fmt.Sprintf(
				"%s: %d hour%s & %d minute%s",
				class,
				playTime.Hours,
				bungie_utils.Plural(playTime.Hours),
				playTime.Minutes,
				bungie_utils.Plural(playTime.Minutes),
			)
		}

		parts = append(parts, formatted)
	}

	message := fmt.Sprintf("%s - %s", request.Gamertag, strings.Join(parts, " | "))
	return &message, nil
}
