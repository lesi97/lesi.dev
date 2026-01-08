package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func Encrypt(plainText string) (string, error) {
	keyB64 := os.Getenv("ENCRYPTION_KEY")
	if keyB64 == "" {
		return "", errors.New("encryption key missing")
	}

	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)

	out := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(out), nil
}