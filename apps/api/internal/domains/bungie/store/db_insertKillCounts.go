package bungie_store

import "context"

type dbKillCounts struct {
	MembershipID 	string
	WeaponID		string
	PVPKills     	*int
	PVEKills     	*int
	TrialsKills  	*int
	WeaponHash 		string
}

func (s *BungieStore) insertKillCounts(dbData *dbKillCounts) {
	query := `
		INSERT INTO destiny_weapon_kill_counts 
			(membership_id, weapon_id, pvp_kills, pve_kills, trials_kills, weapon_hash)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		ON CONFLICT (membership_id, weapon_id) DO UPDATE
		SET 
			pvp_kills = COALESCE(EXCLUDED.pvp_kills, destiny_weapon_kill_counts.pvp_kills),
			pve_kills = COALESCE(EXCLUDED.pve_kills, destiny_weapon_kill_counts.pve_kills),
			trials_kills = COALESCE(EXCLUDED.trials_kills, destiny_weapon_kill_counts.trials_kills),
			weapon_hash = EXCLUDED.weapon_hash
	`
	_, err := s.DB.Exec(context.Background(), query, 
		dbData.MembershipID, 
		dbData.WeaponID, 
		dbData.PVPKills, 
		dbData.PVEKills, 
		dbData.TrialsKills, 
		dbData.WeaponHash,
	)
	if err != nil {
		s.Logger.Printf("insertKillCounts failed: %v", err)
	}
}
