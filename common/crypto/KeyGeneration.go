package crypto

import (
	"crypto/rand"
	"crypto/rsa"
)

const keySize int = 2048

type RsaKeyPair struct {
	VerificationKey *PrivateKey
	BroadcastKey *PublicKey
}

func GenerateKeyPair() (*RsaKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	return &RsaKeyPair{
		VerificationKey: privateKey,
		BroadcastKey: &privateKey.PublicKey,
	}, nil
}
