package utils

import "strings"

func processPerks(rawData map[string]bungieWeaponRaw) []dbPerk {
	var perks []dbPerk

	for _, item := range rawData {
		if item.ItemTypeDisplayName == "" || !strings.Contains(item.ItemTypeDisplayName, "Perk") {
			continue
		}

		perk := dbPerk{
			HashID:      item.Hash,
			Name:        item.DisplayProperties.Name,
			Description: item.DisplayProperties.Description,
			ItemType:    item.ItemTypeDisplayName,
		}
		perks = append(perks, perk)
	}

	return perks
}
