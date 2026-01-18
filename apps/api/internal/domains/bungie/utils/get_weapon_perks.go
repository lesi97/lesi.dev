package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func getWeaponPerks(ctx context.Context, database *db.DB, logger *utils.Logger, perkHashIDs []string) (*filteredPerksResult, error) {
	defer logger.LogExecutionTime("DATABASE CALL: getWeaponPerks", time.Now(), ctx)
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

	rows, err := database.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perks := make([]databasePerk, 0, 16)

	for rows.Next() {
		var perk databasePerk
		if err := rows.Scan(&perk.Name, &perk.ItemType); err != nil {
			return nil, err
		}
		perks = append(perks, perk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	filtered := filterWeaponPerks(perks)
	return &filtered, nil
}
