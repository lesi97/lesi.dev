package utils

import "strconv"

func GetPerkHashIDs(plugHashes []individualSocket) []string {
	var perkHashIDs []string
	for _, plug := range plugHashes {
		perkHashIDs = append(perkHashIDs, strconv.FormatInt(plug.PlugHash, 10))
	}
	return perkHashIDs
}
