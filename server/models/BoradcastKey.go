package models

import (
	"ForgettiServer/db/models"
	"fmt"
	"forgetti-common/crypto"
	"time"

	"github.com/google/uuid"
)

type BoradcastKey struct {
	KeyId uuid.UUID
	Expiration time.Time
	Key *crypto.PublicKey
}

func FromDbModel(model *models.KeyRecord, unprotect func(string) (string, error)) (*BoradcastKey, error) {
	serializedKey, err := unprotect(model.SerializedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unprotect key: %w", err)
	}

	publicKey, err := crypto.DeserializePublicKey(serializedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize public key: %w", err)
	}

	parsedKeyId, err := uuid.Parse(model.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key ID: %w", err)
	}

	return &BoradcastKey{
		KeyId: parsedKeyId,
		Expiration: model.Expiration,
		Key: publicKey,
	}, nil
}
