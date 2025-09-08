package encryption

import (
	"forgetti-common/crypto"
	"forgetti-common/logging"
)

func Encrypt(content []byte, key *Key) ([]byte, error) {
	logger := logging.MakeLogger("encryption.Encrypt")
	logger.Verbose("Starting content encryption (%d bytes)", len(content))

	logger.Verbose("Getting key bytes for encryption")
	keyBytes, err := key.GetBytes()
	if err != nil {
		logger.Error("Failed to get key bytes: %v", err)
		return nil, err
	}
	logger.Verbose("Key bytes obtained (%d bytes)", len(keyBytes))

	logger.Verbose("Performing AES-256 encryption")
	encryptedContent, err := crypto.EncryptAes256(content, keyBytes)
	if err != nil {
		logger.Error("AES-256 encryption failed: %v", err)
		return nil, err
	}
	logger.Info("Successfully encrypted content (%d bytes -> %d bytes)", len(content), len(encryptedContent))
	return encryptedContent, nil
}

func Decrypt(content []byte, key *Key) ([]byte, error) {
	logger := logging.MakeLogger("encryption.Decrypt")
	logger.Verbose("Starting content decryption (%d bytes)", len(content))

	logger.Verbose("Getting key bytes for decryption")
	keyBytes, err := key.GetBytes()
	if err != nil {
		logger.Error("Failed to get key bytes: %v", err)
		return nil, err
	}
	logger.Verbose("Key bytes obtained (%d bytes)", len(keyBytes))

	logger.Verbose("Performing AES-256 decryption")
	decryptedContent, err := crypto.DecryptAes256(content, keyBytes)
	if err != nil {
		logger.Error("AES-256 decryption failed: %v", err)
		return nil, err
	}
	logger.Info("Successfully decrypted content (%d bytes -> %d bytes)", len(content), len(decryptedContent))
	return decryptedContent, nil
}
