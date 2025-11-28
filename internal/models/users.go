package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           uint64         `gorm:"primaryKey;column:id"`
	Email        *string        `gorm:"column:email;type:text;uniqueIndex:ux_users_email;not null"`
	PasswordHash *string        `gorm:"column:password_hash;type:text"`
	Name         *string        `gorm:"column:name;type:text"`
	Metadata     datatypes.JSON `gorm:"column:metadata;type:jsonb"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (User) TableName() string {
	return "users"
}
