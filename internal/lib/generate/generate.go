package generate

import (
	"crypto/rand"
	"encoding/hex"
)

// RandString is random string generation
func RandString(length int) string {
	b := make([]byte, length/2)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
