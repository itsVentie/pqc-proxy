package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateAuthToken(secret, salt string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyAuthToken(token, secret, salt string) bool {
	expected := GenerateAuthToken(secret, salt)
	return hmac.Equal([]byte(token), []byte(expected))
}
