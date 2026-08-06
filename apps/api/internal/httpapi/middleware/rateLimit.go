package middleware

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

type clientRate struct {
	count       int
	windowStart time.Time
	lastSeen    time.Time
}

func RateLimit() func(http.Handler) http.Handler {
	const (
		maxRequests    = 120
		window         = time.Minute
		cleanupEvery   = 5 * time.Minute
		staleClientTTL = 15 * time.Minute
	)

	var mu sync.Mutex
	clients := make(map[string]*clientRate)

	ticker := time.NewTicker(cleanupEvery)
	go func() {
		for range ticker.C {
			mu.Lock()
			for key, client := range clients {
				if time.Since(client.lastSeen) > staleClientTTL {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthcheck" {
				next.ServeHTTP(w, r)
				return
			}

			if shouldBypassRateLimit(r) {
				next.ServeHTTP(w, r)
				return
			}

			now := time.Now()
			clientIP := clientIPFromRequest(r)

			mu.Lock()
			client, ok := clients[clientIP]
			if !ok {
				client = &clientRate{
					count:       0,
					windowStart: now,
					lastSeen:    now,
				}
				clients[clientIP] = client
			}

			if now.Sub(client.windowStart) >= window {
				client.count = 0
				client.windowStart = now
			}

			client.count++
			client.lastSeen = now
			allowed := client.count <= maxRequests
			remaining := max(maxRequests-client.count, 0)

			resetAfter := max(client.windowStart.Add(window).Sub(now), 0)

			resetSeconds := int(resetAfter / time.Second)
			if resetAfter%time.Second != 0 {
				resetSeconds++
			}
			mu.Unlock()

			w.Header().Set("RateLimit-Limit", strconv.Itoa(maxRequests))
			w.Header().Set("RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("RateLimit-Reset", strconv.Itoa(resetSeconds))

			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(resetSeconds))
				utils.TextResponse(w, http.StatusTooManyRequests, "Too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func shouldBypassRateLimit(r *http.Request) bool {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("RATE_LIMIT_DISABLED")), "true") {
		return true
	}

	if os.Getenv("GO_ENV") == "production" {
		return false
	}

	return isLoopbackClientIP(clientIPFromRequest(r))
}

func clientIPFromRequest(r *http.Request) string {
	clientIP := r.RemoteAddr

	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		clientIP = strings.TrimSpace(parts[0])
	}

	realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if realIP != "" {
		clientIP = realIP
	}

	host, _, err := net.SplitHostPort(clientIP)
	if err == nil {
		clientIP = host
	}

	return clientIP
}

func isLoopbackClientIP(clientIP string) bool {
	parsedIP := net.ParseIP(strings.TrimSpace(clientIP))
	return parsedIP != nil && parsedIP.IsLoopback()
}
