package commons

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashBytesToSHA256(input []byte) []byte {
	h := sha256.New()
	return h.Sum(nil)
}

func ConvertBytesToString(input []byte) string {
	return hex.EncodeToString(input)
}
