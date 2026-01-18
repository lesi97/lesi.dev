package utils

type dbWeapons struct {
	ID                  string
	DisplayName         string
	ItemTypeDisplayName string
	FlavorText          string
	BucketTypeHash      string
	TierTypeHash        string
	TierTypeName        string
	TierType            string
	TalentGridHash      string
}

type dbPerk struct {
	Name        string
	Description string
	ItemType    string
	HashID      string
}

type bungieWeaponRaw struct {
	Hash              string `json:"hash"`
	ItemType          int    `json:"itemType"`
	DisplayProperties struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"displayProperties"`
	ItemTypeDisplayName string `json:"itemTypeDisplayName"`
	FlavorText          string `json:"flavorText"`
	Inventory           struct {
		BucketTypeHash string `json:"bucketTypeHash"`
		TierTypeHash   string `json:"tierTypeHash"`
		TierTypeName   string `json:"tierTypeName"`
		TierType       string `json:"tierType"`
	} `json:"inventory"`
	TalentGrid struct {
		TalentGridHash string `json:"talentGridHash"`
	} `json:"talentGrid"`
}
