package crypto

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	pair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}
	if len(pair.PublicKey) == 0 || len(pair.PrivateKey) == 0 {
		t.Error("Keys are empty")
	}
}
