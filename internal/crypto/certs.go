package crypto

import (
	"fmt"

	"github.com/cloudflare/circl/kem/kyber/kyber768"
)

type PQCKeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

func GenerateKeyPair() (*PQCKeyPair, error) {
	pk, sk, err := kyber768.Scheme().GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PQC keys: %w", err)
	}

	pkBytes, _ := pk.MarshalBinary()
	skBytes, _ := sk.MarshalBinary()

	return &PQCKeyPair{
		PublicKey:  pkBytes,
		PrivateKey: skBytes,
	}, nil
}
