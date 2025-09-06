package models

import (
	"forgetti-common/crypto"
	"time"
)

type NewKeyEncryptionResult struct {
	KeyId string
	Expiration time.Time
	VerificationKey *crypto.PrivateKey
	EncryptedContent string
}
