package utils

import (
	"encoding/json"
	"net/http"
)

func MarshalRequestCaptureHeaders(headers http.Header) (string, error) {
	encodedHeaders, err := json.Marshal(RedactRequestCaptureHeaders(headers))
	if err != nil {
		return "", err
	}

	return string(encodedHeaders), nil
}
