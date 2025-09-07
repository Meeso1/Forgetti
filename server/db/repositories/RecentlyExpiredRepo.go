package repositories

import (
	"ForgettiServer/db/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type RecentlyExpiredRepo struct {
	db *gorm.DB
}

func NewRecentlyExpiredRepo(db *gorm.DB) *RecentlyExpiredRepo {
	return &RecentlyExpiredRepo{db: db}
}

func (s *RecentlyExpiredRepo) Create(id string, expiration time.Time) error {
	record := models.RecentlyExpiredRecord{
		Id:         id,
		Expiration: expiration,
	}

	if err := s.db.Create(&record).Error; err != nil {
		return fmt.Errorf("failed to create recently expired record: %w", err)
	}

	return nil
}

func (s *RecentlyExpiredRepo) GetById(id string) (*models.RecentlyExpiredRecord, error) {
	var record models.RecentlyExpiredRecord
	err := s.db.Where("id = ?", id).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get recently expired record: %w", err)
	}

	return &record, nil
}

func (s *RecentlyExpiredRepo) Delete(id string) error {
	result := s.db.Where("id = ?", id).Delete(&models.RecentlyExpiredRecord{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete recently expired record: %w", result.Error)
	}
	return nil
}

func (s *RecentlyExpiredRepo) DeleteBefore(cutoffTime time.Time) error {
	result := s.db.Where("expiration < ?", cutoffTime).Delete(&models.RecentlyExpiredRecord{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete old recently expired records: %w", result.Error)
	}
	return nil
}
