package services

import (
	"ForgettiServer/models"
	"forgetti-common/crypto"
	"time"
	"github.com/google/uuid"
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
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	key := models.BoradcastKey{
		KeyId: uuid.New(),
		Expiration: expiration,
		Key: keyPair.BroadcastKey,
	}

	e.keyStore.StoreKey(key)

	encryptedContent, err := crypto.Encrypt(content, key.Key)
	if err != nil {
		return nil, err
	}

	return &models.NewKeyEncryptionResult{
		KeyId:            key.KeyId.String(),
		Expiration:       key.Expiration,
		VerificationKey:  keyPair.VerificationKey,
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
