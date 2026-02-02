package utils

import (
	"net/http"
	"net/url"
)

func GetRequestOrigin(r *http.Request) (string, bool) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		return origin, true
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		return "", false
	}

	parsed, err := url.Parse(referer)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", false
	}

	return parsed.Scheme + "://" + parsed.Host, true
}
