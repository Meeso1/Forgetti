package services

import (
	"ForgettiServer/models"
	"crypto/rand"
	"crypto/rsa"
	"time"
	"github.com/google/uuid"
)

const keySize int = 2048

type RsaKeyPair struct {
	VerificationKey *rsa.PrivateKey
	BroadcastKey *rsa.PublicKey
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

func GenerateKey(expiration time.Time) (*models.BoradcastKey, *rsa.PrivateKey, error) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	key := models.BoradcastKey{
		KeyId: uuid.New(),
		Expiration: expiration,
		Key: keyPair.BroadcastKey,
	}

	return &key, keyPair.VerificationKey, nil
}
