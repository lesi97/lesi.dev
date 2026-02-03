package middleware

import "strings"

func GetResponseBody(body []byte) string {
	responseBody := strings.TrimSpace(string(body))
	if responseBody == "" {
		return "<empty>"
	}

	return responseBody
}
