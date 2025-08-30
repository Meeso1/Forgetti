package services

import (
	"ForgettiServer/models"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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

// TODO: Remove
func (e *EncryptorImpl) Decrypt(encryptedContent string, serializedVerificationKey string) (string, error) {
	deserializedVerificationKey, err := DeserializePrivateKey(serializedVerificationKey)
	if err != nil {
		return "", err
	}
	
	decryptedContent, err := Decrypt(encryptedContent, deserializedVerificationKey)
	if err != nil {
		return "", err
	}

	return decryptedContent, nil
}

func Encrypt(content string, key *rsa.PublicKey) (string, error) {
	contentBytes := []byte(content)

	// Create a deterministic seed from the content using SHA-256
	hasher := sha256.New()
	hasher.Write(contentBytes)
	seed := hasher.Sum(nil)

	// Use the first 32 bytes of the hash as our deterministic seed
	// This ensures the same content always produces the same ciphertext
	deterministicReader := &DeterministicReader{seed: seed, pos: 0}

	// Encrypt using RSA-OAEP with SHA-256 and our deterministic seed
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), deterministicReader, key, contentBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt content: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func Decrypt(encryptedContent string, key *rsa.PrivateKey) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedContent)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 content: %w", err)
	}

	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, encryptedBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt content: %w", err)
	}

	return string(decryptedBytes), nil
}

// DeterministicReader provides a deterministic source of "randomness" for RSA-OAEP
// by cycling through the provided seed bytes. This ensures deterministic encryption.
type DeterministicReader struct {
	seed []byte
	pos  int
}

func (r *DeterministicReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = r.seed[r.pos%len(r.seed)]
		r.pos++
	}
	return len(p), nil
}
