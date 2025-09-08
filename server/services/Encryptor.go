package services

import (
	"ForgettiServer/models"
	"forgetti-common/crypto"
	"forgetti-common/logging"
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
	logger := logging.MakeLogger("services.CreateEncryptor")
	logger.Verbose("Creating new Encryptor service")
	return &EncryptorImpl{
		keyStore: keyStore,
	}
}

func (e *EncryptorImpl) CreateNewKeyAndEncrypt(content string, expiration time.Time) (*models.NewKeyEncryptionResult, error) {
	logger := logging.MakeLogger("services.Encryptor.CreateNewKeyAndEncrypt")
	logger.Verbose("Creating new key with expiration: %s", expiration.Format("2006-01-02 15:04:05"))

	logger.Verbose("Generating RSA key pair")
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		logger.Error("Failed to generate key pair: %v", err)
		return nil, err
	}
	logger.Verbose("Key pair generated successfully")

	keyId := uuid.New()
	key := models.BoradcastKey{
		KeyId:      keyId,
		Expiration: expiration,
		Key:        keyPair.BroadcastKey,
	}
	logger.Verbose("Generated KeyId: %s", keyId.String())

	logger.Verbose("Storing key in key store")
	err = e.keyStore.StoreKey(key)
	if err != nil {
		logger.Error("Failed to store key: %v", err)
		return nil, err
	}
	logger.Verbose("Key stored successfully")

	logger.Verbose("Encrypting content with RSA key")
	encryptedContent, err := crypto.EncryptRsa(content, key.Key)
	if err != nil {
		logger.Error("Failed to encrypt content: %v", err)
		return nil, err
	}
	logger.Verbose("Content encrypted successfully")

	result := &models.NewKeyEncryptionResult{
		KeyId:            key.KeyId.String(),
		Expiration:       key.Expiration,
		VerificationKey:  keyPair.VerificationKey,
		EncryptedContent: encryptedContent,
	}
	logger.Info("Successfully created new key and encrypted content. KeyId: %s", result.KeyId)
	return result, nil
}

func (e *EncryptorImpl) EncryptWithExistingKey(content string, keyId string) (string, error) {
	logger := logging.MakeLogger("services.Encryptor.EncryptWithExistingKey")
	logger.Verbose("Encrypting with existing KeyId: %s", keyId)

	logger.Verbose("Retrieving key from key store")
	key, err := e.keyStore.GetKey(keyId)
	if err != nil {
		logger.Error("Failed to get key from store: %v", err)
		return "", err
	}
	logger.Verbose("Key retrieved successfully from store")

	logger.Verbose("Encrypting content with existing RSA key")
	encryptedContent, err := crypto.EncryptRsa(content, key.Key)
	if err != nil {
		logger.Error("Failed to encrypt content with existing key: %v", err)
		return "", err
	}
	logger.Verbose("Content encrypted successfully with existing key")
	logger.Info("Successfully encrypted with existing key. KeyId: %s", keyId)

	return encryptedContent, nil
}
