package repositories

import (
	"ForgettiServer/db/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type KeyRepo struct {
	db *gorm.DB
}

func NewKeyRepo(db *gorm.DB) *KeyRepo {
	return &KeyRepo{db: db}
}

func (s *KeyRepo) Create(id string, expiration time.Time, serializedKey string) error {
	var existing models.KeyRecord
	err := s.db.Where("id = ?", id).First(&existing).Error
	if err == nil {
		return fmt.Errorf("key with id %s already exists", id)
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing key: %w", err)
	}

	record := models.KeyRecord{
		Id:            id,
		Expiration:    expiration,
		SerializedKey: serializedKey,
	}

	if err := s.db.Create(&record).Error; err != nil {
		return fmt.Errorf("failed to create key record: %w", err)
	}

	return nil
}

func (s *KeyRepo) GetById(id string) (*models.KeyRecord, error) {
	var record models.KeyRecord
	err := s.db.Where("id = ?", id).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key record: %w", err)
	}

	return &record, nil
}

func (s *KeyRepo) Delete(id string) error {
	result := s.db.Where("id = ?", id).Delete(&models.KeyRecord{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete key record: %w", result.Error)
	}
	return nil
}

func (s *KeyRepo) DeleteExpiredBefore(cutoffTime time.Time) error {
	result := s.db.Where("expiration < ?", cutoffTime).Delete(&models.KeyRecord{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired keys: %w", result.Error)
	}
	return nil
}
