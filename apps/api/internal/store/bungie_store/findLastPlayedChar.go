package bungie_store

func findLastPlayedChar(s *SupabaseBungieStore, chars characters) *character {
	var latest *character
	for _, char := range chars.Data {
		if latest == nil || char.DateLastPlayed.After(latest.DateLastPlayed) {
			copy := char
			latest = &copy
		}

		go func() {
			dbCharData := userCharacters{
				MembershipID:  char.MembershipID,
				CharacterID:   char.CharacterID,
				CharacterType: getCharacterType(char.ClassType),
				MinutesPlayed: char.MinutesPlayedTotal,
			}
			s.insertUserCharacters(&dbCharData)
		}()

	}
	return latest
}
