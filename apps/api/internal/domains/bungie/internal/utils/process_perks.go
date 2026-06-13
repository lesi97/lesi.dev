package utils

import "strings"

func processPerks(rawData map[string]bungieWeaponRaw) []dbPerk {
	var perks []dbPerk

	for hash, item := range rawData {
		itemType := weaponSocketItemType(item)
		if itemType == "" {
			continue
		}

		perk := dbPerk{
			HashID:      item.hashID(hash),
			Name:        item.DisplayProperties.Name,
			Description: item.DisplayProperties.Description,
			ItemType:    itemType,
		}
		perks = append(perks, perk)
	}

	return perks
}

func weaponSocketItemType(item bungieWeaponRaw) string {
	itemType := strings.TrimSpace(item.ItemTypeDisplayName)
	name := strings.TrimSpace(item.DisplayProperties.Name)
	if name == "" {
		return ""
	}

	lowerItemType := strings.ToLower(itemType)
	switch {
	case strings.Contains(lowerItemType, "perk"),
		strings.Contains(lowerItemType, "trait"),
		strings.Contains(lowerItemType, "intrinsic"),
		strings.Contains(lowerItemType, "frame"),
		strings.Contains(lowerItemType, "barrel"),
		strings.Contains(lowerItemType, "magazine"),
		strings.Contains(lowerItemType, "battery"),
		strings.Contains(lowerItemType, "bowstring"),
		strings.Contains(lowerItemType, "arrow"),
		strings.Contains(lowerItemType, "scope"),
		strings.Contains(lowerItemType, "sight"),
		strings.Contains(lowerItemType, "blade"),
		strings.Contains(lowerItemType, "guard"),
		strings.Contains(lowerItemType, "haft"),
		strings.Contains(lowerItemType, "grip"),
		strings.Contains(lowerItemType, "stock"),
		strings.Contains(lowerItemType, "weapon mod"),
		strings.Contains(lowerItemType, "shader"),
		strings.Contains(lowerItemType, "ornament"):
		return itemType
	}

	plugCategory := strings.ToLower(item.Plug.PlugCategoryIdentifier)
	switch {
	case strings.Contains(plugCategory, "skins"):
		return "Weapon Ornament"
	case strings.Contains(plugCategory, "shader"):
		return "Shader"
	}
	if plugCategory != "" {
		if itemType != "" {
			return itemType
		}
		return "Plug"
	}

	return ""
}
