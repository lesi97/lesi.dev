package bungie_store

import (
	"context"
	"fmt"
	"strings"
)

type databasePerk struct {
	Name     string `json:"name"`
	ItemType string `json:"item_type"`
}

type weaponPerks struct {
	Perks []databasePerk `json:"perks"`
	Mods []databasePerk `json:"mods"`
}

type filteredPerksResult struct {
	Perks  []databasePerk
	Mods   []databasePerk
	Shaders []databasePerk
	Ornaments []databasePerk
}



func (supabase *SupabaseBungieStoreStore) getWeaponPerks(ctx context.Context, perkHashIDs []string) (*filteredPerksResult, error) {
	if len(perkHashIDs) == 0 {
		return nil, fmt.Errorf("perk list not provided")
	}

	placeholders := make([]string, len(perkHashIDs))
	args := make([]interface{}, len(perkHashIDs))
	for i, id := range perkHashIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`
		select 
			name, 
			item_type 
		from destiny_weapon_perks 
		where hash_id in (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := supabase.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := &weaponPerks{
		Perks: []databasePerk{},
	}

	for rows.Next() {
		var perk databasePerk
		if err := rows.Scan(&perk.Name, &perk.ItemType); err != nil {
			return nil, err
		}
		data.Perks = append(data.Perks, perk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	filtered := filterWeaponPerks(data.Perks)
	return &filtered, nil
}

func filterWeaponPerks(perks []databasePerk) filteredPerksResult  {
	var filtered []databasePerk
	var mods []databasePerk
	var shaders []databasePerk
	var ornaments []databasePerk

	for _, perk := range perks {
		name := strings.ToLower(perk.Name)
		itemType := strings.ToLower(perk.ItemType)

		switch {
		case strings.Contains(itemType, "ornament"), strings.Contains(itemType, "Ornament") :
			ornaments = append(ornaments, perk)

		case strings.Contains(itemType, "shader"):
			shaders = append(shaders, perk)

		case strings.Contains(itemType, "mod"):
			mods = append(mods, perk)

		case strings.Contains(name, "masterwork"),
			strings.Contains(name, "catalyst"),
			strings.Contains(name, "default shader"),
			strings.Contains(name, "unknown"),
			strings.Contains(name, "socket"),
			strings.Contains(name, "shaped weapon"),
			strings.Contains(name, "crucible tracker"),
			strings.Contains(name, "kill tracker"),
			strings.Contains(name, "trials memento"):
			// Ignored
		default:
			filtered = append(filtered, perk)
		}
	}

	return filteredPerksResult{
		Perks:  filtered,
		Mods:   mods,
		Shaders: shaders,
		Ornaments: ornaments,
	}
}