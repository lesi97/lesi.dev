package bungie_store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type bungieSearchResponse struct {
	IconPath                    string `json:"iconPath"`
	CrossSaveOverride           int    `json:"crossSaveOverride"`
	ApplicableMembershipTypes   []int  `json:"applicableMembershipTypes"`
	IsPublic                    bool   `json:"isPublic"`
	MembershipType              int    `json:"membershipType"`
	MembershipID                string `json:"membershipId"`
	DisplayName                 string `json:"displayName"`
	BungieGlobalDisplayName     string `json:"bungieGlobalDisplayName"`
	BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
}

type bungieSearch struct {
	Response 		[]bungieSearchResponse  `json:"Response"`
	ErrorCode       int    					`json:"ErrorCode"`
	ThrottleSeconds int   					`json:"ThrottleSeconds"`
	ErrorStatus     string 					`json:"ErrorStatus"`
	Message         string 					`json:"Message"`
	MessageData     struct {} 				`json:"MessageData"`

}

func (s *BungieStore) getUserFromBungieByGamertag(ctx context.Context, id string, platform string) (*bungieSearch, error) {
	platformEnum := switchPlatforms(platform)
	cacheKey := fmt.Sprintf("bungie:searchdestinyplayer:%s:%s", id, platformEnum)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		result := &bungieSearch{}
		if err := json.Unmarshal([]byte(cached), result); err == nil {
			s.Logger.Printf("CACHE HIT getUserFromBungieByGamertag %s", cacheKey)
			utils.PrintPrettyJSON(result)
			return result, nil
		}
		_ = s.redis.Del(ctx, cacheKey).Err()
	} else {
		if err != redis.Nil {
			return nil, err
		}
	}

	escapedID := url.PathEscape(id)
	url := fmt.Sprintf("%s/Platform/Destiny2/SearchDestinyPlayer/%s/%s/", s.url, platformEnum, escapedID)

	body, err := s.bungieGET(url)
	if err != nil {
		fmt.Println("ERROR in getUserFromBungieByGamertag")
		return nil, err
	}

	result := &bungieSearch{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	b, err := json.Marshal(result)
	if err == nil {
		_ = s.redis.Set(ctx, cacheKey, b, 168 * time.Hour).Err()
	}

	go func() {
		if len(result.Response) > 0 {
			user := bungieDBData{
				BungieID: id,
				MembershipID: result.Response[0].MembershipID,
				PreferredPlatform: int64(result.Response[0].MembershipType),
				FriendlyName: result.Response[0].BungieGlobalDisplayName,
			}
		s.insertDestinyUser(&user)
		}
	}()

	return result, nil
}
