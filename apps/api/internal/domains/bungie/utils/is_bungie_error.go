package utils

import (
	"encoding/json"
	"fmt"
)

func isBungieError(body []byte) (bool, error) {
	var response errorResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return false, nil
	}
	if response.ErrorCode != 1 {
		return true, fmt.Errorf("bungie error: %s (%d)", response.Message, response.ErrorCode)
	}
	return false, nil
}
