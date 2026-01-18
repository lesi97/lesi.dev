package twitch_frontend

import (
	"crypto/sha256"
	"encoding/base64"
)

func PKCEChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
