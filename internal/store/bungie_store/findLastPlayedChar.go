package bungie_store

func findLastPlayedChar(chars characters) *character {
	var latest *character
	for _, char := range chars.Data {
		if latest == nil || char.DateLastPlayed.After(latest.DateLastPlayed) {
			copy := char
			latest = &copy
		}
	}
	return latest
}
