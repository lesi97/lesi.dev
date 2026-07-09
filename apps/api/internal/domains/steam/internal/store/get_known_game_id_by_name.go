package store

func GetKnownGameIDByName(normalisedGameName string) (string, bool) {
	knownGameIDs := map[string]string{
		"rocket league": "252950",

		"resident evil 9": "3764200",
		"re9":             "3764200",

		"marathon beta":    "4254230",
		"deadlock":         "1422450",
		"ghost of tushima": "2215430",

		"call of duty vanguard": "1985820",
		"cod vanguard":          "1985820",

		"rainbow six siege": "359550",
		"rainbow six":       "359550",
		"r6":                "359550",

		"the first descendent": "2074920",
		"first descendent":     "2074920",

		"hello kitty island adventures": "2495100",

		"ac black flag":                        "3751950",
		"ac resynced":                          "3751950",
		"ac resyncd":                           "3751950",
		"black flag":                           "3751950",
		"assassins creed black flag resynced":  "3751950",
		"assassin's creed black flag resynced": "3751950",
		"assassins creed black flag resyncd":   "3751950",
		"assassins creed black flag":           "3751950",
		"assassins creed resynced":             "3751950",
		"assassins creed resyncd":              "3751950",
	}

	gameID, ok := knownGameIDs[normalisedGameName]
	if !ok {
		return "", false
	}

	return gameID, true
}
