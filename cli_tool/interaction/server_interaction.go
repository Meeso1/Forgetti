package interaction

import (
	"forgetti-common/crypto"
	"Forgetti/encryption"
	"Forgetti/models"
	"time"
	"fmt"
)

type KeyGenerationResult struct {
	EncryptedKeyHash string
	Metadata models.Metadata
}

func GenerateKeyAndEncrypt(serverAddress string, key string, expiration time.Time) (*KeyGenerationResult, error) {
	remoteClient := NewRemoteClient(serverAddress)

	keyHash, err := encryption.HashRemotePartForEncryption(key)
	if err != nil {
		return nil, err
	}

	response, err := remoteClient.NewKey(keyHash, expiration)
	if err != nil {
		return nil, err
	}

	if err := validateEncryptedKeyHash(keyHash, response.EncryptedContent, response.Metadata.VerificationKey); err != nil {
		return nil, err
	}

	return &KeyGenerationResult{
		EncryptedKeyHash: response.EncryptedContent,
		Metadata: models.Metadata{
			KeyId: response.Metadata.KeyId,
			Expiration: response.Metadata.Expiration,
			VerificationKey: response.Metadata.VerificationKey,
			ServerAddress: serverAddress,
		},
	}, nil
}

func EncryptWithExistingKey(serverAddress string, key string, metadata *models.Metadata) (string, error) {
	remoteClient := NewRemoteClient(serverAddress)

	keyHash, err := encryption.HashRemotePartForEncryption(key)
	if err != nil {
		return "", err
	}

	response, err := remoteClient.Encrypt(keyHash, metadata.KeyId)
	if err != nil {
		return "", err
	}

	if err := validateEncryptedKeyHash(keyHash, response.EncryptedContent, metadata.VerificationKey); err != nil {
		return "", err
	}

	return response.EncryptedContent, nil
}

func validateEncryptedKeyHash(keyHash string, encrypted string, serializedKey string) error {
	verificationKey, err := crypto.DeserializePrivateKey(serializedKey)
	if err != nil {
		return err
	}

	decryptedKeyHash, err := crypto.DecryptRsa(encrypted, verificationKey)
	if err != nil {
		return err
	}

	if decryptedKeyHash != keyHash {
		return fmt.Errorf("decrypted key hash does not match the original key hash")
	}

	return nil
}