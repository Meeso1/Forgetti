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
	key []byte
}

func NewDataProtection(config *config.Config) (DataProtection, error) {
	key, err := base64.StdEncoding.DecodeString(config.DataProtection.KeyBase64)
	if err != nil {
		return nil, err
	}

	return &DataProtectionImpl{key: key}, nil
}

func (d *DataProtectionImpl) Protect(data string) (string, error) {
	encrypted, err := crypto.EncryptAes256([]byte(data), d.key)
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

	decrypted, err := crypto.DecryptAes256(encrypted, d.key)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
