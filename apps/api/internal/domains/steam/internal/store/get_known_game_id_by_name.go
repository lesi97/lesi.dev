package store

func GetKnownGameIDByName(normalisedGameName string) (string, bool) {
	knownGameIDs := map[string]string{
		"rocket league": "252950",
	}

	gameID, ok := knownGameIDs[normalisedGameName]
	if !ok {
		return "", false
	}

	return gameID, true
}
