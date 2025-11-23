package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type WalletTransaction struct {
	ID            uint64          `gorm:"primaryKey"`
	Chain         string          `gorm:"type:text;default:ethereum;not null"`
	TxHash        string          `gorm:"type:text;not null;index:ux_chain_txhash,unique"`
	BlockNumber   *int64          `gorm:"type:bigint"`
	BlockTime     *time.Time      `gorm:"column:block_timestamp;type:timestamptz"`
	FromAddress   *string         `gorm:"type:text"`
	ToAddress     *string         `gorm:"type:text"`
	Value         decimal.Decimal `gorm:"type:numeric(50,18)"`
	GasUsed       decimal.Decimal `gorm:"type:numeric(30,0)"`
	GasPrice      decimal.Decimal `gorm:"type:numeric(30,0)"`
	Method        *string         `gorm:"type:text"`
	Decoded       datatypes.JSON  `gorm:"type:jsonb"`
	Status        *int16          `gorm:"type:smallint"`
	RawTx         datatypes.JSON  `gorm:"type:jsonb"`
	WalletAddress *string         `gorm:"type:text;index"`
	CreatedAt     time.Time       `gorm:"type:timestamptz;default:now()"`
}

func (WalletTransaction) TableName() string { return "wallet_transactions" }
