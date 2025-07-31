package bungie_store

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateID(id string) bool {
	idRegex := regexp.MustCompile(`^[\w !@#$%^&*()_+={}\[\]:;"'<>,.?/\\|-]+#[0-9]{4}$`)
	invalidRegex := regexp.MustCompile(`\b(drop|alter|delete|insert|update|create|select|truncate|exec|union)\b`)
	if invalidRegex.MatchString(strings.ToLower(id)) {
		return false
	}
	if !idRegex.MatchString(id) {
		return false
	}
	return true
}

// interface is the arugment type to accept null values and return -1
func getPlatformEnum(platform interface{}) string {
	platformString, ok := platform.(string)
	if !ok {
		return "-1"
	}

	switch strings.ToLower(platformString) {
	case "xb", "xbox":
		return "1"
	case "ps", "playstation":
		return "2"
	case "pc":
		return "3"
	case "bnet":
		return "4"
	case "st", "steam":
		return "5"
	case "demon":
		return "10"
	default:
		return "-1"
	}
}


type PlayTime struct {
	Minutes	int
	Hours 	int
}
func formatPlayTime(playTime int) PlayTime {
	hours := playTime / 60
	minutes := playTime % 60
	return PlayTime{
		Hours:   hours,
		Minutes: minutes,
	}
}

func isBungieError(body []byte) (bool, error) {
	var resposne errorResponse
	err := json.Unmarshal(body, &resposne) 
	if err != nil {
		return false, nil
	}
	if resposne.ErrorCode != 1 {
		return true, fmt.Errorf("bungie error: %s (%d)", resposne.Message, resposne.ErrorCode)
	}
	return false, nil
}


func getCharacterType(characterType int) string {
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

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
