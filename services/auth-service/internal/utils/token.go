package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateVerificationToken create a pair token: Raw Token (send through mail) và Hashed Token (stored in DB)
func GenerateVerificationToken() (rawToken string, hashedToken string, err error) {
	// 1. Create random 32 bytes
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	// 2. Encode to string
	rawToken = hex.EncodeToString(randomBytes)

	// 3. Hash rawToken usingg SHA-256
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken = hex.EncodeToString(hash[:])

	return rawToken, hashedToken, nil
}
