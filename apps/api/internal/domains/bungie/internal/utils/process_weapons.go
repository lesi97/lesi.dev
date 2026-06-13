package utils

func processWeapons(rawData map[string]bungieWeaponRaw) []dbWeapons {
	var weapons []dbWeapons

	for hash, item := range rawData {
		if item.ItemType != 3 {
			continue
		}

		weapon := dbWeapons{
			ID:                  item.hashID(hash),
			DisplayName:         item.DisplayProperties.Name,
			ItemTypeDisplayName: item.ItemTypeDisplayName,
			FlavorText:          item.FlavorText,
			BucketTypeHash:      item.Inventory.BucketTypeHash.String(),
			TierTypeHash:        item.Inventory.TierTypeHash.String(),
			TierTypeName:        item.Inventory.TierTypeName,
			TierType:            item.Inventory.TierType.Int(),
			TalentGridHash:      item.TalentGrid.TalentGridHash.String(),
		}
		weapons = append(weapons, weapon)
	}

	return weapons
}
