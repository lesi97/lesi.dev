package utils

import (
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func FindLastPlayedCharacter(database *db.DB, logger *utils.Logger, chars characters) *character {
	var latest *character
	for _, char := range chars.Data {
		if latest == nil || char.DateLastPlayed.After(latest.DateLastPlayed) {
			copy := char
			latest = &copy
		}

		go func(ch character) {
			InsertUserCharacters(
				database,
				logger,
				ch.MembershipID,
				ch.CharacterID,
				GetCharacterType(ch.ClassType),
				ch.MinutesPlayedTotal,
			)
		}(char)
	}
	return latest
}
