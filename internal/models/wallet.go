package models

import (
	"time"

	"gorm.io/datatypes"
)

type Wallet struct {
	Address       string         `gorm:"primaryKey;type:text"`
	Chain         string         `gorm:"type:text;default:ethereum;not null"`
	ChainCoins    int64          `gorm:"type:bigint;default:0"`
	LastIndexedAt time.Time      `gorm:"type:timestamptz;default:now()"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;default:now()"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;default:now()"`
}

func (Wallet) TableName() string { return "wallets" }
