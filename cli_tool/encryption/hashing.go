package encryption

import (
	"forgetti-common/crypto"
	"encoding/base64"
	"fmt"
)

const beforeEncryptionSalt string = "before_encryption"
const remoteSalt string = "remote"
const localSalt string = "local"

func HashRemotePartForEncryption(key string, version string) (string, error) {
	if version != "1" {
		return "", fmt.Errorf("unsupported version: %s", version)
	}

	bytes, err := crypto.HashToSize(key, beforeEncryptionSalt, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HashEncryptedRemotePart(key string, version string) ([]byte, error) {
	if version != "1" {
		return nil, fmt.Errorf("unsupported version: %s", version)
	}

	return crypto.HashToSize(key, remoteSalt, 16)
}

func HashLocalPart(key string, version string) ([]byte, error) {
	if version != "1" {
		return nil, fmt.Errorf("unsupported version: %s", version)
	}

	return crypto.HashToSize(key, localSalt, 16)
}


