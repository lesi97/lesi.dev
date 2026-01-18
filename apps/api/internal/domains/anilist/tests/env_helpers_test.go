package tests

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load(".env.local")
}

func envOrFail(t *testing.T, key string) string {
	t.Helper()

	value := os.Getenv(key)
	if value == "" {
		t.Skipf("missing env var %s", key)
	}
	t.Setenv(key, value)
	return value
}
