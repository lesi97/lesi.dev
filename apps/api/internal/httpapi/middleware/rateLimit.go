package middleware

import (
	"net"
	"net/http"
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

			now := time.Now()
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
			mu.Unlock()

			if !allowed {
				utils.TextResponse(w, http.StatusTooManyRequests, "Too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
