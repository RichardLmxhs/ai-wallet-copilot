package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           int64          `gorm:"primaryKey;column:id"`
	Email        *string        `gorm:"column:email"`
	PasswordHash *string        `gorm:"column:password_hash"`
	Name         *string        `gorm:"column:name"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	Metadata     datatypes.JSON `gorm:"column:metadata"`
}

func (User) TableName() string {
	return "users"
}
