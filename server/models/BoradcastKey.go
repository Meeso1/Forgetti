package models

import (
	"crypto/rsa"
	"time"
	"github.com/google/uuid"
)

type BoradcastKey struct {
	KeyId uuid.UUID
	Expiration time.Time
	Key *rsa.PublicKey
}
