package bungie_store

import (
	"context"
	"fmt"
)

type user struct {
	MembershipID   string
	MembershipType int
	DisplayName    string
	Source         string
}

type result struct {
	user *user
	err  error
}

func (s *SupabaseBungieStoreStore) getUser(ctx context.Context, gt string) (*user, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan result, 2)

	go func() {
		dbUser, err := s.getUserFromDatabaseByGamertag(ctx, gt)
		if err != nil || dbUser == nil {
			ch <- result{nil, err}
			return
		}
		user := &user{
			MembershipID:   dbUser.MembershipID,
			MembershipType: int(dbUser.PreferredPlatform),
			DisplayName:    dbUser.FriendlyName,
			Source:         "db",
		}
		ch <- result{user, nil}
	}()

	go func() {
		apiUser, err := getUserFromBungieByGamertag(gt)
		if err != nil || apiUser == nil || len(apiUser.Response) == 0 {
			ch <- result{nil, err}
			return
		}
		bungie := apiUser.Response[0]
		user := &user{
			MembershipID:   bungie.MembershipID,
			MembershipType: bungie.MembershipType,
			DisplayName:    bungie.BungieGlobalDisplayName,
			Source:         "api",
		}
		ch <- result{user, nil}
	}()

	for range 2 {
		select {
		case result := <-ch:
			if result.err == nil && result.user != nil {
				cancel()
				return result.user, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("no result found from either source")
}