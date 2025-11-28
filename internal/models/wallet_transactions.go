package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type WalletTransaction struct {
	ID             uint64          `gorm:"primaryKey"`
	Chain          string          `gorm:"type:text;default:ethereum;not null"`
	TxHash         string          `gorm:"column:tx_hash;type:text;not null;index:ux_wallet_transactions_chain_txhash,unique"`
	BlockNumber    *int64          `gorm:"column:block_number;type:bigint;index:idx_wallet_transactions_block"`
	BlockTimestamp *time.Time      `gorm:"type:timestamptz"`
	FromAddress    *string         `gorm:"column:from_address;type:text"`
	ToAddress      *string         `gorm:"column:to_address;type:text"`
	Value          decimal.Decimal `gorm:"type:numeric(50,18)"`
	GasUsed        decimal.Decimal `gorm:"column:gas_used;type:numeric(30,0)"`
	GasPrice       decimal.Decimal `gorm:"column:gas_price;type:numeric(30,0)"`
	Method         *string         `gorm:"type:text"`
	Decoded        datatypes.JSON  `gorm:"type:jsonb"`
	Status         *int16          `gorm:"type:smallint"`
	RawTx          datatypes.JSON  `gorm:"column:raw_tx;type:jsonb"`
	WalletAddress  *string         `gorm:"column:wallet_address;type:text;index:idx_wallet_transactions_wallet"`
	CreatedAt      time.Time       `gorm:"type:timestamptz;default:now()"`
}

func (WalletTransaction) TableName() string { return "wallet_transactions" }
