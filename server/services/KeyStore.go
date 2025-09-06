package services

import (
	"ForgettiServer/config"
	"ForgettiServer/errors"
	"ForgettiServer/models"
	"fmt"
	"time"
)

type KeyStore interface {
	StoreKey(key models.BoradcastKey) error
	GetKey(keyId string) (*models.BoradcastKey, error)
}

// TODO: Store keys in a database
type KeyStoreImpl struct {
	keys                    map[string]models.BoradcastKey
	recentlyExpired         map[string]time.Time
	recentlyExpiredDuration time.Duration
}

func CreateKeyStore(cfg *config.Config) KeyStore {
	return &KeyStoreImpl{
		keys:                    make(map[string]models.BoradcastKey),
		recentlyExpired:         make(map[string]time.Time),
		recentlyExpiredDuration: time.Duration(cfg.KeyStore.RecentlyExpiredDurationHours) * time.Hour,
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
			if expiration.Before(time.Now().Add(-k.recentlyExpiredDuration)) {
				delete(k.recentlyExpired, keyId)
				return nil, errors.KeyNotFoundError(keyId)
			}

			return nil, errors.KeyExpiredError(keyId, expiration)
		}

		return nil, errors.KeyNotFoundError(keyId)
	}

	// Expiration check
	if key.Expiration.Before(time.Now()) {
		// If key expired more than the configured duration ago, remove and return not found
		if key.Expiration.Before(time.Now().Add(-k.recentlyExpiredDuration)) {
			delete(k.keys, keyId)
			return nil, errors.KeyNotFoundError(keyId)
		}

		k.recentlyExpired[keyId] = key.Expiration
		delete(k.keys, keyId)

		return nil, errors.KeyExpiredError(keyId, key.Expiration)
	}

	return &key, nil
}
