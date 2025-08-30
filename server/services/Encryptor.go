package services

import (
	"ForgettiServer/models"
	"crypto/rsa"
	// "crypto/sha256"
	// "crypto/rand"
	"time"
)

type Encryptor interface {
	CreateNewKeyAndEncrypt(content string, expiration time.Time) (*models.NewKeyEncryptionResult, error)
	EncryptWithExistingKey(content string, keyId string) (string, error)
}

type EncryptorImpl struct {
	keyStore KeyStore
}

func CreateEncryptor(keyStore KeyStore) Encryptor {
	return &EncryptorImpl{
		keyStore: keyStore,
	}
}

func (e *EncryptorImpl) CreateNewKeyAndEncrypt(content string, expiration time.Time) (*models.NewKeyEncryptionResult, error) {
	key, verificationKey, err := GenerateKey(expiration)
	if err != nil {
		return nil, err
	}

	e.keyStore.StoreKey(*key)

	encryptedContent, err := Encrypt(content, key.Key)
	if err != nil {
		return nil, err
	}

	return &models.NewKeyEncryptionResult{
		KeyId: key.KeyId.String(),
		Expiration: key.Expiration,
		VerificationKey: verificationKey,
		EncryptedContent: encryptedContent,
	}, nil
}

func (e *EncryptorImpl) EncryptWithExistingKey(content string, keyId string) (string, error) {
	key, err := e.keyStore.GetKey(keyId)
	if err != nil {
		return "", err
	}
	
	encryptedContent, err := Encrypt(content, key.Key)
	if err != nil {
		return "", err
	}

	return encryptedContent, nil
}

func Encrypt(content string, key *rsa.PublicKey) (string, error) {
	// TODO: Implement

	return content, nil
}