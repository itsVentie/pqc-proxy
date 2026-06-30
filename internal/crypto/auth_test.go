package crypto

import "testing"

func TestAuthLogic(t *testing.T) {
	secret := "super-secret-password"
	salt := "random-salt-123"

	token := GenerateAuthToken(secret, salt)

	if !VerifyAuthToken(token, secret, salt) {
		t.Error("Token verification failed for correct key")
	}

	if VerifyAuthToken(token, "wrong-password", salt) {
		t.Error("Token verification should have failed for wrong key")
	}
}
