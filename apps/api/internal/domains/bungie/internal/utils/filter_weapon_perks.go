package utils

import "strings"

func filterWeaponPerks(perks []databasePerk) filteredPerksResult {
	var filtered []databasePerk
	var mods []databasePerk
	var shaders []databasePerk
	var ornaments []databasePerk

	for _, perk := range perks {
		name := strings.ToLower(perk.Name)
		itemType := strings.ToLower(perk.ItemType)

		switch {
		case strings.Contains(itemType, "ornament"), strings.Contains(itemType, "Ornament"):
			ornaments = append(ornaments, perk)

		case strings.Contains(itemType, "shader"):
			shaders = append(shaders, perk)

		case strings.Contains(itemType, "mod"):
			mods = append(mods, perk)

		case strings.Contains(name, "masterwork"),
			strings.Contains(name, "catalyst"),
			strings.Contains(name, "default shader"),
			strings.Contains(name, "gear tier upgrade"),
			strings.Contains(name, "unknown"),
			strings.Contains(name, "socket"),
			strings.Contains(name, "shaped weapon"),
			strings.Contains(name, "crucible tracker"),
			strings.Contains(name, "kill tracker"),
			strings.Contains(name, "trials memento"),
			strings.Contains(itemType, "combat flair"),
			itemType == "plug":
		default:
			filtered = append(filtered, perk)
		}
	}

	return filteredPerksResult{
		Perks:     filtered,
		Mods:      mods,
		Shaders:   shaders,
		Ornaments: ornaments,
	}
}
