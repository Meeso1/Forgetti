package services

import (
	"ForgettiServer/config"
	"ForgettiServer/db/repositories"
	dbModels "ForgettiServer/db/models"
	"ForgettiServer/errors"
	"ForgettiServer/models"
	"fmt"
	"forgetti-common/crypto"
	"time"
)

type KeyStore interface {
	StoreKey(key models.BoradcastKey) error
	GetKey(keyId string) (*models.BoradcastKey, error)
}

type KeyStoreImpl struct {
	keyRepo                 *repositories.KeyRepo
	recentlyExpiredRepo     *repositories.RecentlyExpiredRepo
	dataProtection          DataProtection
	recentlyExpiredDuration time.Duration
}

func NewKeyStore(
	keyRepo *repositories.KeyRepo,
	recentlyExpiredRepo *repositories.RecentlyExpiredRepo,
	dataProtection DataProtection,
	cfg *config.Config,
) KeyStore {
	return &KeyStoreImpl{
		keyRepo:                 keyRepo,
		recentlyExpiredRepo:     recentlyExpiredRepo,
		dataProtection:          dataProtection,
		recentlyExpiredDuration: time.Duration(cfg.KeyStore.RecentlyExpiredDurationHours) * time.Hour,
	}
}

func (k *KeyStoreImpl) StoreKey(key models.BoradcastKey) error {
	serializedKey, err := crypto.SerializePublicKey(key.Key)
	if err != nil {
		return fmt.Errorf("failed to serialize key: %w", err)
	}

	protectedKey, err := k.dataProtection.Protect(serializedKey)
	if err != nil {
		return fmt.Errorf("failed to protect key: %w", err)
	}

	return k.keyRepo.Create(key.KeyId.String(), key.Expiration, protectedKey)
}

func (k *KeyStoreImpl) GetKey(keyId string) (*models.BoradcastKey, error) {
	record, err := k.keyRepo.GetById(keyId)
	if err != nil {
		return nil, fmt.Errorf("failed to get key from database: %w", err)
	}

	if record != nil {
		record, err = k.checkExpiration(record)
		if err != nil {
			return nil, err
		}

		result, err := models.FromDbModel(record, k.dataProtection.Unprotect)
		if err != nil {
			return nil, fmt.Errorf("failed to convert database model to broadcast key: %w", err)
		}

		return result, nil
	}

	// Check recently expired
	expiredRecord, err := k.recentlyExpiredRepo.GetById(keyId)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently expired record: %w", err)
	}
	if expiredRecord != nil {
		return nil, errors.KeyExpiredError(keyId, expiredRecord.Expiration)
	}
	return nil, errors.KeyNotFoundError(keyId)
}

func (k *KeyStoreImpl) checkExpiration(record *dbModels.KeyRecord) (*dbModels.KeyRecord, error) {
	if !record.Expiration.Before(time.Now()) {
		return record, nil
	}

	// Expired too long ago, treat as not found
	if record.Expiration.Before(time.Now().Add(-k.recentlyExpiredDuration)) {
		if err := k.keyRepo.Delete(record.Id); err != nil {
			return nil, fmt.Errorf("failed to delete expired key: %w", err)
		}
		return nil, errors.KeyNotFoundError(record.Id)
	}

	// Move to recently expired
	if err := k.recentlyExpiredRepo.Create(record.Id, record.Expiration); err != nil {
		return nil, fmt.Errorf("failed to create recently expired record: %w", err)
	}
	if err := k.keyRepo.Delete(record.Id); err != nil {
		return nil, fmt.Errorf("failed to delete expired key: %w", err)
	}

	return nil, errors.KeyExpiredError(record.Id, record.Expiration)
}

func (k *KeyStoreImpl) CleanupExpiredKeys() error {
	cutoffTime := time.Now().Add(-k.recentlyExpiredDuration)
	if err := k.recentlyExpiredRepo.DeleteBefore(cutoffTime); err != nil {
		return fmt.Errorf("failed to cleanup recently expired records: %w", err)
	}

	if err := k.keyRepo.DeleteExpiredBefore(cutoffTime); err != nil {
		return fmt.Errorf("failed to cleanup old expired keys: %w", err)
	}

	return nil
}
