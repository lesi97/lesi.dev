package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type userResult struct {
	user *User
	err  error
}

func GetUser(
	ctx context.Context,
	database *db.DB,
	logger *utils.Logger,
	redis *redis.Client,
	httpClient *http.Client,
	baseURL string,
	clientID string,
	gt string,
	platform string,
) (*User, error) {
	defer logger.LogExecutionTime("MATRIX: getUser", time.Now(), ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan userResult, 2)

	go func() {
		dbUser, err := getUserFromDatabaseByGamertag(ctx, database, logger, redis, gt)
		if err != nil || dbUser == nil {
			if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil {
				logger.Printf("ERROR in Matrix - getUserFromDatabase: %v\n", err)
			}
			ch <- userResult{nil, fmt.Errorf("no users found in db with gt: %v", gt)}
			return
		}
		user := &User{
			MembershipID:   dbUser.MembershipID,
			MembershipType: int(dbUser.PreferredPlatform),
			DisplayName:    dbUser.FriendlyName,
			Source:         "db",
		}
		ch <- userResult{user, nil}
	}()

	go func() {
		apiUser, err := getUserFromBungieByGamertag(ctx, database, logger, redis, httpClient, baseURL, clientID, gt, platform)
		if err != nil || apiUser == nil || apiUser.Response == nil || len(apiUser.Response) == 0 {
			if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil {
				logger.Printf("ERROR in Matrix - getUserFromBungie: %v\n", err)
			}
			ch <- userResult{nil, fmt.Errorf("bungie response was empty")}
			return
		}
		bungie := apiUser.Response[0]
		user := &User{
			MembershipID:   bungie.MembershipID,
			MembershipType: bungie.MembershipType,
			DisplayName:    bungie.BungieGlobalDisplayName,
			Source:         "api",
		}
		ch <- userResult{user, nil}
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
