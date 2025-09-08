package interaction

import (
	"Forgetti/encryption"
	"Forgetti/models"
	"fmt"
	"forgetti-common/crypto"
	"forgetti-common/logging"
	"time"
)

type KeyGenerationResult struct {
	EncryptedKeyHash string
	Metadata         models.Metadata
}

func GenerateKeyAndEncrypt(serverAddress string, key string, expiration time.Time) (*KeyGenerationResult, error) {
	logger := logging.MakeLogger("server_interaction.GenerateKeyAndEncrypt")
	remoteClient := NewRemoteClient(serverAddress)

	logger.Verbose("Hashing key for server interaction with pre-remote hash algorithm")
	keyHash, err := encryption.HashRemotePartForEncryption(key, models.CurrentAlgVersion().PreRemoteHash)
	if err != nil {
		logger.Error("Failed to hash key: %v", err)
		return nil, err
	}
	logger.Verbose("Key hashed successfully")

	logger.Verbose("Making new key request to server %s with expiration %s", serverAddress, expiration.Format("2006-01-02 15:04:05"))
	response, err := remoteClient.NewKey(keyHash, expiration)
	if err != nil {
		logger.Error("Failed to create new key on server: %v", err)
		return nil, err
	}
	logger.Verbose("New key created on server with KeyId: %s", response.Metadata.KeyId)

	logger.Verbose("Validating encrypted key hash")
	if err := validateEncryptedKeyHash(keyHash, response.EncryptedContent, response.Metadata.VerificationKey); err != nil {
		logger.Error("Key hash validation failed: %v", err)
		return nil, err
	}
	logger.Verbose("Key hash validation successful")

	result := &KeyGenerationResult{
		EncryptedKeyHash: response.EncryptedContent,
		Metadata:         models.ToFileMetadata(response.Metadata, serverAddress),
	}
	logger.Info("Successfully generated and encrypted key. KeyId: %s, Expiration: %s", result.Metadata.KeyId, result.Metadata.Expiration.Format("2006-01-02 15:04:05"))
	return result, nil
}

func EncryptWithExistingKey(serverAddress string, key string, metadata *models.Metadata) (string, error) {
	logger := logging.MakeLogger("server_interaction.EncryptWithExistingKey")
	remoteClient := NewRemoteClient(serverAddress)
	versions := models.ParseAlgVersion(metadata.AlgVersion)

	logger.Verbose("Hashing key for existing key encryption with KeyId: %s", metadata.KeyId)
	keyHash, err := encryption.HashRemotePartForEncryption(key, versions.PreRemoteHash)
	if err != nil {
		logger.Error("Failed to hash key for existing key encryption: %v", err)
		return "", err
	}
	logger.Verbose("Key hashed successfully for existing key")

	logger.Verbose("Making encrypt request to server %s for KeyId: %s", serverAddress, metadata.KeyId)
	response, err := remoteClient.Encrypt(keyHash, metadata.KeyId)
	if err != nil {
		logger.Error("Failed to encrypt with existing key on server: %v", err)
		return "", err
	}
	logger.Verbose("Encrypt request successful")

	logger.Verbose("Validating encrypted key hash for existing key")
	if err := validateEncryptedKeyHash(keyHash, response.EncryptedContent, metadata.VerificationKey); err != nil {
		logger.Error("Key hash validation failed for existing key: %v", err)
		return "", err
	}
	logger.Verbose("Key hash validation successful for existing key")
	logger.Info("Successfully encrypted with existing key. KeyId: %s", metadata.KeyId)

	return response.EncryptedContent, nil
}

func validateEncryptedKeyHash(keyHash string, encrypted string, serializedKey string) error {
	logger := logging.MakeLogger("server_interaction.validateEncryptedKeyHash")

	logger.Verbose("Deserializing verification key")
	verificationKey, err := crypto.DeserializePrivateKey(serializedKey)
	if err != nil {
		logger.Error("Failed to deserialize verification key: %v", err)
		return err
	}
	logger.Verbose("Verification key deserialized successfully")

	logger.Verbose("Decrypting key hash with RSA")
	decryptedKeyHash, err := crypto.DecryptRsa(encrypted, verificationKey)
	if err != nil {
		logger.Error("Failed to decrypt key hash: %v", err)
		return err
	}
	logger.Verbose("Key hash decrypted successfully")

	logger.Verbose("Validating decrypted key hash matches original")
	if decryptedKeyHash != keyHash {
		logger.Error("Key hash validation failed - decrypted hash does not match original")
		return fmt.Errorf("decrypted key hash does not match the original key hash")
	}
	logger.Verbose("Key hash validation passed")

	return nil
}
