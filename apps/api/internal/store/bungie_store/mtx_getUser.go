package bungie_store

import (
	"context"
	"fmt"
	"time"
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

func (s *BungieStore) getUser(ctx context.Context, gt string, platform string) (*user, error) {
	defer s.Logger.LogExecutionTime("MATRIX: getUser", time.Now(), ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan result, 2)

	go func() {
		dbUser, err := s.getUserFromDatabaseByGamertag(ctx, gt)
		if err != nil || dbUser == nil {
			s.Logger.Printf("ERROR in Matrix - getUserFromDatabase: %v\n", err)
			ch <- result{nil, fmt.Errorf("no users found in db with gt: %v", gt)}
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
		apiUser, err := s.getUserFromBungieByGamertag(ctx, gt, platform)
		if err != nil || apiUser == nil || apiUser.Response == nil || len(apiUser.Response) == 0 {
			s.Logger.Printf("ERROR in Matrix - getUserFromBungie: %v\n", err)
			ch <- result{nil, fmt.Errorf("bungie response was empty")}
			return
		} else {			
			bungie := apiUser.Response[0]
			user := &user{
				MembershipID:   bungie.MembershipID,
				MembershipType: bungie.MembershipType,
				DisplayName:    bungie.BungieGlobalDisplayName,
				Source:         "api",
			}
			ch <- result{user, nil}
		}
	}()

	for range 2 {
		select {
		case res := <-ch:
			if res.err == nil && res.user != nil {
				cancel()
				return res.user, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("unable to find player %v", gt)
}