package store

func GetKnownGameIDByName(normalisedGameName string) (string, bool) {
	knownGameIDs := map[string]string{
		"rocket league": "252950",
		"resident evil 9": "3764200",
		"re9": "3764200",
		"marathon beta": "4254230",
	}

	gameID, ok := knownGameIDs[normalisedGameName]
	if !ok {
		return "", false
	}

	return gameID, true
}
