package bungie_store

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func getBungieProfileByMembershipID(membershipID string, preferredPlatform string, components string) (*BungieProfile, error) {
	url := fmt.Sprintf("%s/Platform/Destiny2/%s/Profile/%s/?components=%s", bungie_url, preferredPlatform, membershipID, components)

	body, err := bungieGET(url)
	if err != nil {
		return nil, err
	}

	result := &BungieProfile{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	return result, nil
}