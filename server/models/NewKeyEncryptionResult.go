package models

import (
	"crypto/rsa"
	"time"
)

type NewKeyEncryptionResult struct {
	KeyId string
	Expiration time.Time
	VerificationKey *rsa.PrivateKey
	EncryptedContent string
}
