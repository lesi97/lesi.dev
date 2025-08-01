package bungie_store

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
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
func getPlatformEnum(platform interface{}, user *user) string {
	platformString, ok := platform.(string)
	if !ok {
		return strconv.Itoa(user.MembershipType)
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
		return strconv.Itoa(user.MembershipType)
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

func getPerkHashIDs(plugHashes []individualSocket) []string {
	var perkHashIDs []string
	for _, plug := range plugHashes {
		perkHashIDs = append(perkHashIDs, strconv.FormatInt(plug.PlugHash, 10))
	}
	return perkHashIDs
}


func getPerkNames(perks []databasePerk) string {
	var perkNames []string
	for _, perk := range perks {
		if strings.Contains(perk.ItemType, "Enhanced") {
			perkName := fmt.Sprintf("Enhanced %s", perk.Name)
			perkNames = append(perkNames, perkName)
		} else {
			perkNames = append(perkNames, perk.Name)
		}
	}
	return strings.Join(perkNames, ", ")
}

func generateString(gt string, weapon *weaponResult, category string, killCount int) string {
	var responseMessage string;

	perkNamesString := getPerkNames(weapon.weaponPerks.Perks)

	responseMessage += fmt.Sprintf("%s: ", gt)
	responseMessage += fmt.Sprintf("%s | ", weapon.weaponData.DisplayName)
	responseMessage += fmt.Sprintf("Perks: %s ", perkNamesString)

	if len(weapon.weaponPerks.Mods) != 0 {
		responseMessage += fmt.Sprintf("| Mod: %s ", weapon.weaponPerks.Mods[0].Name)
	}

	if len(weapon.weaponPerks.Shaders) != 0 {
		responseMessage += fmt.Sprintf("| Shader: %s ", weapon.weaponPerks.Shaders[0].Name)
	}

	if len(weapon.weaponPerks.Ornaments) != 0 {
		responseMessage += fmt.Sprintf("| Ornament: %s ", weapon.weaponPerks.Ornaments[0].Name)
	}

	if category != ""  {
		responseMessage += fmt.Sprintf("| %s Kill Count: %v", category, humanize.Comma(int64(killCount)),)
	}

	return responseMessage
}