package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString returns a cryptographically random alphanumeric string of length n.
func RandomString(n int) string {
	result := make([]byte, n)
	for i := range result {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphanumeric))))
		result[i] = alphanumeric[idx.Int64()]
	}
	return string(result)
}

// RandomID returns a UUID-like random hex string.
func RandomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:])
}
