package bungie_store

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type dbWeapons struct {
	ID 					string
	DisplayName 		string
	ItemTypeDisplayName string
	FlavorText 			string
	BucketTypeHash 		string
	TierTypeHash 		string
	TierTypeName 		string
	TierType 			string
	TalentGridHash 		string
}

type dbPerk struct {
	Name 		string
	Description string
	ItemType 	string
	HashID 		string
}

type bungieWeaponRaw struct {
	Hash               string `json:"hash"`
	ItemType           int    `json:"itemType"`
	DisplayProperties  struct {
		Name string `json:"name"`
		Description string `json:"description"`
	} `json:"displayProperties"`
	ItemTypeDisplayName string `json:"itemTypeDisplayName"`
	FlavorText           string `json:"flavorText"`
	Inventory            struct {
		BucketTypeHash string `json:"bucketTypeHash"`
		TierTypeHash   string `json:"tierTypeHash"`
		TierTypeName   string `json:"tierTypeName"`
		TierType       string `json:"tierType"`
	} `json:"inventory"`
	TalentGrid struct {
		TalentGridHash string `json:"talentGridHash"`
	} `json:"talentGrid"`
}

func (supabase *SupabaseBungieStore) getNewWeapons() {
	urlPath, err := getManifestURL()
	if err != nil {
		fmt.Printf("failed to generate manifest URL")
		return
	}

	url := fmt.Sprintf("%s%s", bungie_url, *urlPath)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to get manifest")
		return
	}

	var rawData map[string]bungieWeaponRaw
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		fmt.Printf("failed to decode data")
		return
	}

	weapons := processWeapons(rawData)
	perks := processPerks(rawData)
	supabase.insertWeapons(weapons)
	supabase.insertPerks(perks)
}

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


func (s *SupabaseBungieStore) insertWeapons(weapons []dbWeapons) {
	for _, weapon := range weapons {
		query := `
			INSERT INTO destiny_weapons 
				(id, display_name, item_type_display_name, flavor_text, bucket_type_hash, tier_type_hash, tier_type_name, tier_type, talent_grid_hash)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id)
			DO NOTHING
		`
		_, err := s.db.Exec(context.Background(), query,
			weapon.ID, 
			weapon.DisplayName, 
			weapon.ItemTypeDisplayName, 
			weapon.FlavorText,
			weapon.BucketTypeHash, 
			weapon.TierTypeHash, 
			weapon.TierTypeName,
			weapon.TierType, 
			weapon.TalentGridHash,
		)
		if err != nil {
			s.logger.Printf("insertWeapons failed for ID %s: %v\n", weapon.ID, err)
		}
	}
}

func (s *SupabaseBungieStore) insertPerks(perks []dbPerk) {
	for _, perk := range perks {
		query := `
			INSERT INTO destiny_weapon_perks
				(name, description, item_type, hash_id)
			VALUES 
				($1, $2, $3, $4)
			ON CONFLICT (hash_id)
			DO NOTHING
		`
		_, err := s.db.Exec(context.Background(), query, perk.Name, perk.Description, perk.ItemType, perk.HashID)
		if err != nil {
			s.logger.Printf("insertPerks failed for %s: %v\n", perk.Name, err)
		}
	}
}
