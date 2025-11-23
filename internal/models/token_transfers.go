package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type TokenTransfer struct {
	ID           uint64          `gorm:"primaryKey"`
	Chain        string          `gorm:"type:text;default:ethereum;not null"`
	TxHash       string          `gorm:"type:text;not null"`
	LogIndex     int             `gorm:"type:int"`
	TokenAddress *string         `gorm:"type:text"`
	TokenType    *string         `gorm:"type:text"`
	FromAddress  *string         `gorm:"type:text;index"`
	ToAddress    *string         `gorm:"type:text;index"`
	Amount       decimal.Decimal `gorm:"type:numeric(50,18)"`
	TokenID      *string         `gorm:"type:text"`
	Decimals     int             `gorm:"type:int"`
	Metadata     datatypes.JSON  `gorm:"type:jsonb"`
	CreatedAt    time.Time       `gorm:"type:timestamptz;default:now()"`
}

func (TokenTransfer) TableName() string { return "token_transfers" }
