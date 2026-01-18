package utils

import "strings"

func switchPlatforms(platform string) string {
	switch strings.ToLower(platform) {
	case "xb", "xbox":
		return "1"
	case "ps", "playstation":
		return "2"
	case "pc":
		return "3"
	case "bnet", "b-net", "b_net", "battlenet", "battle-net", "battle_net":
		return "4"
	case "st", "steam":
		return "5"
	case "demon":
		return "10"
	default:
		return "-1"
	}
}
