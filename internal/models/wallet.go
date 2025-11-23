package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Wallet struct {
	Address       string          `gorm:"primaryKey;type:text"`
	Chain         string          `gorm:"type:text;default:ethereum;not null"`
	LastBlock     *int64          `gorm:"type:bigint"`
	TotalValueUSD decimal.Decimal `gorm:"type:numeric(30,8)"`
	TokenCount    int             `gorm:"type:int;default:0"`
	NFTCount      int             `gorm:"type:int;default:0"`
	LastIndexedAt time.Time       `gorm:"type:timestamptz;default:now()"`
	Metadata      datatypes.JSON  `gorm:"type:jsonb"`
	CreatedAt     time.Time       `gorm:"type:timestamptz;default:now()"`
	UpdatedAt     time.Time       `gorm:"type:timestamptz;default:now()"`
}

func (Wallet) TableName() string { return "wallets" }
