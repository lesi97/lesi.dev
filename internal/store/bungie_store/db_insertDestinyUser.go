package bungie_store

import "context"

func (supabase *SupabaseBungieStore) insertDestinyUser(user *bungieDBData) {
	query := `
		INSERT INTO destiny_users 
			(membership_id, bungie_id, preferred_platform, friendly_name)
		VALUES 
			($1, $2, $3, $4)
		ON CONFLICT (bungie_id) DO NOTHING
	`
	_, err := supabase.db.Exec(context.Background(), query, user.MembershipID, user.BungieID, user.PreferredPlatform, user.FriendlyName)
	if err != nil {
		supabase.logger.Printf("insertUser failed: %v", err)
	}
}
