package bungie_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
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

func getUserFromBungieByGamertag(id string) (*bungieSearch, error) {
	escapedID := url.PathEscape(id)
	url := fmt.Sprintf("%s/Destiny2/SearchDestinyPlayer/-1/%s/", bungie_url, escapedID)

	body, err := bungieGET(url)
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

	return result, nil
}
