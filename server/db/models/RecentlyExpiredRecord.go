package models

import "time"

type RecentlyExpiredRecord struct {
	Id         string    `gorm:"primarykey;column:id" json:"id"`
	Expiration time.Time `gorm:"column:expiration;not null" json:"expiration"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (RecentlyExpiredRecord) TableName() string {
	return "recently_expired"
}

func init() {
	RegisterModel(&RecentlyExpiredRecord{})
}
