package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(content []byte, key Key) ([]byte, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, content, nil), nil
}

func Decrypt(content []byte, key Key) ([]byte, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(content) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := content[:nonceSize], content[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func getKeyBytes(key Key) ([]byte, error) {
	keyBytes, err := key.GetBytes()
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("function returned %d bytes, expected 32 bytes", len(keyBytes))
	}

	return keyBytes, nil
}
