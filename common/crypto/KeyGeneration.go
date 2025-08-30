package crypto

import (
	"crypto/rand"
	"crypto/rsa"
)

const keySize int = 2048

type KeyPair struct {
	VerificationKey *PrivateKey
	BroadcastKey *PublicKey
}

func GenerateKeyPair() (*KeyPair, error) {
	keyPair, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		VerificationKey: keyPair,
		BroadcastKey: &keyPair.PublicKey,
	}, nil
}
