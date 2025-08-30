package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Encrypt encrypts content using RSA-OAEP with deterministic behavior
func Encrypt(content string, key *PublicKey) (string, error) {
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

// Decrypt decrypts content using RSA-OAEP
func Decrypt(encryptedContent string, key *PrivateKey) (string, error) {
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
