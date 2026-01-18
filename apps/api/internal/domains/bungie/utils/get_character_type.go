package utils

func GetCharacterType(characterType int) string {
	switch characterType {
	case 0:
		return "Titan"
	case 1:
		return "Hunter"
	case 2:
		return "Warlock"
	default:
		return "UNK"
	}
}
