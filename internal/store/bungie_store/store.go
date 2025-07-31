package bungie_store

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lesi97/api.lesi.dev/internal/database"
)

const bungie_url = "https://www.bungie.net/Platform"

type BungieStore interface {
	GetCharacterPlayTime(ctx context.Context, gt string, platform string) (*string, error)
	GetEquippedPrimary(ctx context.Context, gt string, platform string) (*string, error)
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

func (s *SupabaseBungieStoreStore) GetCharacterPlayTime(ctx context.Context, gt string, platform string) (*string, error) {
	preferredPlatform := getPlatformEnum(platform)
	user, err := s.getUser(ctx, gt)
	if err != nil {
		s.logger.Printf("ERROR: getCharacterPlayTime %v\n", err)
		return nil, err
	}
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

	message := fmt.Sprintf("%s - %s", gt, strings.Join(parts, " | "))
	return &message, nil

}

func (s *SupabaseBungieStoreStore) GetEquippedPrimary(ctx context.Context, gt string, platform string) (*string, error) {
	user, err := s.getUser(ctx, gt)
	if err != nil {
		s.logger.Printf("ERROR: getCharacterPlayTime %v\n", err)
		return nil, err
	}
	// user, err := getBungieProfileByMembershipID(id, "3", "200")
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return nil, err
	}
	fmt.Printf("User: %+v\n", user)
	return nil, nil
}