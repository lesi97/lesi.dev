package cache

import (
	"crypto/sha1"
	"encoding/hex"
)

func key(prefix string, raw string) string {
	sum := sha1.Sum([]byte(raw))
	return prefix + ":" + hex.EncodeToString(sum[:])
}

func (c *Cache) Key(prefix string, raw string) string {
	return key(prefix, raw)
}