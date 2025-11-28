package models

import (
	"time"

	"gorm.io/datatypes"
)

type RiskFlag struct {
	ID            uint64         `gorm:"primaryKey;column:id"`
	WalletAddress string         `gorm:"column:wallet_address;type:text;not null;index:idx_risk_flags_wallet"`
	Chain         string         `gorm:"column:chain;type:text;default:ethereum;not null"`
	FlagType      *string        `gorm:"column:flag_type;type:text;not null;index:idx_risk_flags_type"`
	Score         *int           `gorm:"column:score;type:int"`
	Evidence      datatypes.JSON `gorm:"column:evidence;type:jsonb"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamptz;default:now()"`
}

func (RiskFlag) TableName() string {
	return "risk_flags"
}
