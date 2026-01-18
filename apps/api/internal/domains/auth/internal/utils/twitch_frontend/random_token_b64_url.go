package twitch_frontend

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomTokenB64URL(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
