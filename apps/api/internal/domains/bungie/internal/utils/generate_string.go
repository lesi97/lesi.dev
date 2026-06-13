package utils

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func GenerateString(gt string, weapon *WeaponResult, gearTier int, category string, killCount int) string {
	var responseMessage string

	perkNamesString := getPerkNames(weapon.weaponPerks.Perks)

	responseMessage += fmt.Sprintf("%s: ", gt)
	responseMessage += fmt.Sprintf("%s | ", weapon.weaponData.DisplayName)
	if gearTier > 0 {
		responseMessage += fmt.Sprintf("Tier %d | ", gearTier)
	}
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

	if category != "" {
		responseMessage += fmt.Sprintf("| %s Kill Count: %v", category, humanize.Comma(int64(killCount)))
	}

	return responseMessage
}
