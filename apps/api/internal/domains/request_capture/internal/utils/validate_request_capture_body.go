package utils

import "encoding/json"

func ValidateRequestCaptureBody(body string) error {
	var payload json.RawMessage

	return json.Unmarshal([]byte(body), &payload)
}
