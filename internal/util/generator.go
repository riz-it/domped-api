package util

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"github.com/google/uuid"
)

// GenerateUUID generates a random UUID.
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateRandomCode generates a random numeric code with the given length.
func GenerateRandomCode(length int64) string {
	const charset = "0123456789"

	code := make([]byte, length)
	for i := range code {
		randomIndex, err := randomInt(len(charset))
		if err != nil {
			panic("failed to generate random index")
		}
		code[i] = charset[randomIndex]
	}
	return string(code)
}

// GenerateRandomString generates a random alphanumeric string with the given length.
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomString := make([]byte, length)
	for i := range randomString {
		randomIndex, err := randomInt(len(charset))
		if err != nil {
			panic("failed to generate random index")
		}
		randomString[i] = charset[randomIndex]
	}
	return string(randomString)
}

// GenerateRandomHex generates a random hexadecimal string with the given length.
func GenerateRandomHex(length int) string {
	randomBytes := make([]byte, length/2)
	if _, err := rand.Read(randomBytes); err != nil {
		panic("failed to generate random bytes")
	}
	return hex.EncodeToString(randomBytes)
}

// randomInt generates a random integer between 0 and max using crypto/rand.
func randomInt(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}
