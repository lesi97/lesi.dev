package bungie_store

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type manifestDefinitions struct {
	Response struct {
		JsonWorldComponentContentPaths struct {
			En struct {
				DestinyInventoryItemDefinition string
			} `json:"en"`
		} `json:"jsonWorldComponentContentPaths"`
	} `json:"Response"`
}

func (store *SupabaseBungieStore) getManifestURL() (*string, error) {
	url := fmt.Sprintf("%s/Platform/Destiny2/Manifest", store.url)

	body, err := bungieGET(url)
	if err != nil {
		return nil, err
	}

	result := &manifestDefinitions{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	return &result.Response.JsonWorldComponentContentPaths.En.DestinyInventoryItemDefinition, nil
}