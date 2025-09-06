package models

import (
	"time"
	"forgetti-common/crypto"
	"github.com/google/uuid"
)

type BoradcastKey struct {
	KeyId uuid.UUID
	Expiration time.Time
	Key *crypto.PublicKey
}
