package services

import (
	"ForgettiServer/models"
	"forgetti-common/crypto"
	"time"
)

type Encryptor interface {
	CreateNewKeyAndEncrypt(content string, expiration time.Time) (*models.NewKeyEncryptionResult, error)
	EncryptWithExistingKey(content string, keyId string) (string, error)
	Decrypt(encryptedContent string, keyId string) (string, error)
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

	encryptedContent, err := crypto.Encrypt(content, key.Key)
	if err != nil {
		return nil, err
	}

	return &models.NewKeyEncryptionResult{
		KeyId:            key.KeyId.String(),
		Expiration:       key.Expiration,
		VerificationKey:  verificationKey,
		EncryptedContent: encryptedContent,
	}, nil
}

func (e *EncryptorImpl) EncryptWithExistingKey(content string, keyId string) (string, error) {
	key, err := e.keyStore.GetKey(keyId)
	if err != nil {
		return "", err
	}

	encryptedContent, err := crypto.Encrypt(content, key.Key)
	if err != nil {
		return "", err
	}

	return encryptedContent, nil
}

// TODO: Remove
func (e *EncryptorImpl) Decrypt(encryptedContent string, serializedVerificationKey string) (string, error) {
	deserializedVerificationKey, err := DeserializePrivateKey(serializedVerificationKey)
	if err != nil {
		return "", err
	}

	decryptedContent, err := crypto.Decrypt(encryptedContent, deserializedVerificationKey)
	if err != nil {
		return "", err
	}

	return decryptedContent, nil
}
