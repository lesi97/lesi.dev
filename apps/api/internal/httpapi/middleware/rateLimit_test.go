package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitBypassesLocalRequestsOutsideProduction(t *testing.T) {
	t.Setenv("GO_ENV", "development")
	t.Setenv("RATE_LIMIT_DISABLED", "")

	handler := RateLimit()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 125; i++ {
		req := httptest.NewRequest(http.MethodGet, "/v1/scrobbles", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}
}

func TestRateLimitCanBeDisabledWithEnv(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	t.Setenv("RATE_LIMIT_DISABLED", "true")

	handler := RateLimit()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 125; i++ {
		req := httptest.NewRequest(http.MethodGet, "/v1/scrobbles", nil)
		req.RemoteAddr = "203.0.113.10:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}
}

func TestRateLimitLimitsExternalProductionRequests(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	t.Setenv("RATE_LIMIT_DISABLED", "false")

	handler := RateLimit()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 120; i++ {
		req := httptest.NewRequest(http.MethodGet, "/v1/scrobbles", nil)
		req.RemoteAddr = "203.0.113.10:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/scrobbles", nil)
	req.RemoteAddr = "203.0.113.10:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
}
