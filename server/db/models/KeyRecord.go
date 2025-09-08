package models

import "time"

type KeyRecord struct {
	Id            string    `gorm:"primarykey;column:id" json:"id"`
	Expiration    time.Time `gorm:"column:expiration;not null" json:"expiration"`
	SerializedKey string    `gorm:"column:serialized_key;not null" json:"serialized_key"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (KeyRecord) TableName() string {
	return "keys"
}

func init() {
	RegisterModel(&KeyRecord{})
}
