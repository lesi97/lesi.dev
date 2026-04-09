package store

func GetKnownGameIDByName(normalisedGameName string) (string, bool) {
	knownGameIDs := map[string]string{
		"rocket league":         "252950",
		"resident evil 9":       "3764200",
		"re9":                   "3764200",
		"marathon beta":         "4254230",
		"deadlock":              "1422450",
		"ghost of tushima":      "2215430",
		"call of duty vanguard": "1985820",
		"cod vanguard":          "1985820",
		"rainbow six siege":     "359550",
		"rainbow six":           "359550",
		"r6":                    "359550",
	}

	gameID, ok := knownGameIDs[normalisedGameName]
	if !ok {
		return "", false
	}

	return gameID, true
}
