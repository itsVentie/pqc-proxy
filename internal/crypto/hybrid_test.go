package crypto_test

import (
	"bytes"
	"testing"

	"pqc-proxy/internal/crypto"
)

func TestHybridKeyExchange(t *testing.T) {
	clientEcdhPriv, clientMlkemPriv, clientBlob, err := crypto.GenerateClientInception()
	if err != nil {
		t.Fatalf("Failed to generate client inception: %v", err)
	}

	serverMasterKey, serverBlob, err := crypto.ServerHandleInception(clientBlob)
	if err != nil {
		t.Fatalf("Server failed to handle inception: %v", err)
	}

	clientMasterKey, err := crypto.ClientHandleResponse(clientEcdhPriv, clientMlkemPriv, serverBlob)
	if err != nil {
		t.Fatalf("Client failed to handle server response: %v", err)
	}

	if !bytes.Equal(clientMasterKey, serverMasterKey) {
		t.Error("Critical crypto failure: Derived master keys do not match!")
	}

	t.Logf("Success! Shared master key generated (32 bytes): %x", clientMasterKey)
}
