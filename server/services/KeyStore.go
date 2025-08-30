package services

import (
	"ForgettiServer/models"
	"fmt"
	"time"
)

const recentlyExpiredDuration time.Duration = 24 * time.Hour

type KeyStore interface {
	StoreKey(key models.BoradcastKey) error
	GetKey(keyId string) (*models.BoradcastKey, error)
}

// TODO: Store keys in a database
type KeyStoreImpl struct {
	keys map[string]models.BoradcastKey
	recentlyExpired map[string]time.Time
}

func CreateKeyStore() KeyStore {
	return &KeyStoreImpl{
		keys: make(map[string]models.BoradcastKey),
		recentlyExpired: make(map[string]time.Time),
	}
}

func (k *KeyStoreImpl) StoreKey(key models.BoradcastKey) error {
	if _, ok := k.keys[key.KeyId.String()]; ok {
		return fmt.Errorf("key with id %s already exists", key.KeyId.String())
	}

	k.keys[key.KeyId.String()] = key
	return nil
}

func (k *KeyStoreImpl) GetKey(keyId string) (*models.BoradcastKey, error) {
	key, ok := k.keys[keyId]
	if !ok {
		if expiration, ok := k.recentlyExpired[keyId]; ok {
			// Cleanup from recently expired
			if expiration.Before(time.Now().Add(-recentlyExpiredDuration)) {
				delete(k.recentlyExpired, keyId)
				return nil, fmt.Errorf("key with id %s not found", keyId)
			}

			return nil, fmt.Errorf("key with id %s has expired at %s", keyId, expiration.Format(time.RFC3339))
		}

		return nil, fmt.Errorf("key with id %s not found", keyId)
	}

	// Expiration check
	if key.Expiration.Before(time.Now()) {
		// If key expired more than 24 hours ago, remove and return not found
		if key.Expiration.Before(time.Now().Add(-recentlyExpiredDuration)) {
			delete(k.keys, keyId)
			return nil, fmt.Errorf("key with id %s not found", keyId)
		}

		k.recentlyExpired[keyId] = key.Expiration
		delete(k.keys, keyId)

		return nil, fmt.Errorf("key with id %s has expired at %s", keyId, key.Expiration.Format(time.RFC3339))
	}

	return &key, nil
}
