package services

import (
	"ForgettiServer/config"
	"encoding/base64"
	"forgetti-common/crypto"
)

type DataProtection interface {
	Protect(data string) (string, error)
	Unprotect(data string) (string, error)
}

type DataProtectionImpl struct {
	key string
}

func NewDataProtection(config *config.Config) DataProtection {
	return &DataProtectionImpl{key: config.DataProtection.Key}
}

func (d *DataProtectionImpl) Protect(data string) (string, error) {
	keyHash, err := crypto.HashToSize(d.key, "data_protection", 32)
	if err != nil {
		return "", err
	}

	encrypted, err := crypto.EncryptAes256([]byte(data), keyHash)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (d *DataProtectionImpl) Unprotect(data string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	keyHash, err := crypto.HashToSize(d.key, "data_protection", 32)
	if err != nil {
		return "", err
	}

	decrypted, err := crypto.DecryptAes256(encrypted, keyHash)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
