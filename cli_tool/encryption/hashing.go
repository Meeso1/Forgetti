package encryption

import (
	"crypto/sha256"
	"fmt"
	"encoding/base64"
)

// TODO: Things like this should be versioned
const beforeEncryptionSalt string = "before_encryption"
const remoteSalt string = "remote"
const localSalt string = "local"

func HashRemotePartForEncryption(key string) (string, error) {
	bytes, err := hashToSize(key, beforeEncryptionSalt, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HashEncryptedRemotePart(key string) ([]byte, error) {
	return hashToSize(key, remoteSalt, 16)
}

func HashLocalPart(key string) ([]byte, error) {
	return hashToSize(key, localSalt, 16)
}

func hashToSize(key string, salt string, size int) ([]byte, error) {
	if size > 32 {
		return nil, fmt.Errorf("size must be less than or equal to 32: got %d", size)
	}

	hash := sha256.New()
	hash.Write([]byte(key))
	hash.Write([]byte(salt))
	return hash.Sum(nil)[:size], nil
}
