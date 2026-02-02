package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
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

func GetManifestURL(logger *utils.Logger, httpClient *http.Client, clientID string, baseURL string) (*string, error) {
	url := fmt.Sprintf("%s/Platform/Destiny2/Manifest", baseURL)

	body, err := BungieGET(context.Background(), logger, httpClient, clientID, url)
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
