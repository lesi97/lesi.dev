package handler

import (
	"net"
	"net/http"
	"strings"
)

func isLocalAuthRequest(r *http.Request) bool {
	return isLocalHost(r.RemoteAddr)
}

func isLocalHost(value string) bool {
	host := strings.TrimSpace(value)
	if host == "" {
		return false
	}

	if splitHost, _, err := net.SplitHostPort(host); err == nil {
		host = splitHost
	}

	host = strings.Trim(host, "[]")
	if strings.EqualFold(host, "localhost") {
		return true
	}

	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
