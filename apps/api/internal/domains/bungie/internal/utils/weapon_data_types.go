package utils

type weaponData struct {
	HashID       string `json:"hash_id"`
	DisplayName  string `json:"display_name"`
	TierTypeName string `json:"tier_type_name"`
}

type databasePerk struct {
	Name     string `json:"name"`
	ItemType string `json:"item_type"`
}

type filteredPerksResult struct {
	Perks     []databasePerk
	Mods      []databasePerk
	Shaders   []databasePerk
	Ornaments []databasePerk
}

type WeaponResult struct {
	weaponData  *weaponData
	weaponPerks *filteredPerksResult
	err         error
}
