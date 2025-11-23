package models

import (
	"time"

	"gorm.io/datatypes"
)

type RiskFlag struct {
	ID            int64          `gorm:"primaryKey;column:id"`
	WalletAddress string         `gorm:"column:wallet_address"`
	Chain         string         `gorm:"column:chain"`
	FlagType      *string        `gorm:"column:flag_type"`
	Score         *int           `gorm:"column:score"`
	Evidence      datatypes.JSON `gorm:"column:evidence"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
}

func (RiskFlag) TableName() string {
	return "risk_flags"
}
