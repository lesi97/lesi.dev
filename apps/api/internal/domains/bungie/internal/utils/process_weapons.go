package utils

func processWeapons(rawData map[string]bungieWeaponRaw) []dbWeapons {
	var weapons []dbWeapons

	for _, item := range rawData {
		if item.ItemType != 3 {
			continue
		}

		weapon := dbWeapons{
			ID:                  item.Hash,
			DisplayName:         item.DisplayProperties.Name,
			ItemTypeDisplayName: item.ItemTypeDisplayName,
			FlavorText:          item.FlavorText,
			BucketTypeHash:      item.Inventory.BucketTypeHash,
			TierTypeHash:        item.Inventory.TierTypeHash,
			TierTypeName:        item.Inventory.TierTypeName,
			TierType:            item.Inventory.TierType,
			TalentGridHash:      item.TalentGrid.TalentGridHash,
		}
		weapons = append(weapons, weapon)
	}

	return weapons
}
