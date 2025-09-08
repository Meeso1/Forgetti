package encryption

import (
	"forgetti-common/crypto"
	"encoding/base64"
)

// TODO: Things like this should be versioned
const beforeEncryptionSalt string = "before_encryption"
const remoteSalt string = "remote"
const localSalt string = "local"

func HashRemotePartForEncryption(key string) (string, error) {
	bytes, err := crypto.HashToSize(key, beforeEncryptionSalt, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HashEncryptedRemotePart(key string) ([]byte, error) {
	return crypto.HashToSize(key, remoteSalt, 16)
}

func HashLocalPart(key string) ([]byte, error) {
	return crypto.HashToSize(key, localSalt, 16)
}


