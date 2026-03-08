package httpapi

import (
	"net/url"
	"regexp"
	"strings"
)

func RedactSensitiveQueryValues(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		redactedURL := rawURL
		redactedURL = regexp.MustCompile(`(?i)(key=)([^&\s]+)`).ReplaceAllString(redactedURL, `${1}[REDACTED]`)
		redactedURL = regexp.MustCompile(`(?i)(api_key=)([^&\s]+)`).ReplaceAllString(redactedURL, `${1}[REDACTED]`)
		redactedURL = regexp.MustCompile(`(?i)(apikey=)([^&\s]+)`).ReplaceAllString(redactedURL, `${1}[REDACTED]`)
		redactedURL = regexp.MustCompile(`(?i)(token=)([^&\s]+)`).ReplaceAllString(redactedURL, `${1}[REDACTED]`)
		redactedURL = regexp.MustCompile(`(?i)(secret=)([^&\s]+)`).ReplaceAllString(redactedURL, `${1}[REDACTED]`)
		redactedURL = regexp.MustCompile(`(?i)(client_secret=)([^&\s]+)`).ReplaceAllString(redactedURL, `${1}[REDACTED]`)
		return redactedURL
	}

	query := parsedURL.Query()
	for key := range query {
		lowerKey := strings.ToLower(key)
		if lowerKey == "key" || lowerKey == "api_key" || lowerKey == "apikey" || lowerKey == "token" || lowerKey == "secret" || lowerKey == "client_secret" {
			query.Set(key, "[REDACTED]")
		}
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}
