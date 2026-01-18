package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    fmt.Println(exPath)
	err = godotenv.Load(".env.local")
	fmt.Printf("err: %v\n", err)
}

func envOrFail(t *testing.T, key string) string {
	t.Helper()

	value := os.Getenv(key)
	if value == "" {
		t.Fatalf("missing env var %s", key)
	}
	t.Setenv(key, value)
	return value
}
