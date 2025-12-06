package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", nil
	}
	return hex.EncodeToString(b), nil
}