package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type TokenTransfer struct {
	ID           uint64          `gorm:"primaryKey"`
	Chain        string          `gorm:"type:text;default:ethereum;not null"`
	TxHash       string          `gorm:"column:tx_hash;type:text;not null;index:ux_token_transfers_txhash_logindex,unique:idx_txhash_logindex"`
	LogIndex     int             `gorm:"column:log_index;type:int;uniqueIndex:idx_txhash_logindex"`
	TokenAddress *string         `gorm:"column:token_address;type:text;index:idx_token_transfers_token"`
	TokenType    *string         `gorm:"column:token_type;type:text"`
	FromAddress  *string         `gorm:"column:from_address;type:text;index:idx_token_transfers_from"`
	ToAddress    *string         `gorm:"column:to_address;type:text;index:idx_token_transfers_to"`
	Amount       decimal.Decimal `gorm:"type:numeric(50,18)"`
	TokenID      *string         `gorm:"column:token_id;type:text"`
	Decimals     int             `gorm:"type:int"`
	Metadata     datatypes.JSON  `gorm:"type:jsonb"`
	CreatedAt    time.Time       `gorm:"type:timestamptz;default:now()"`
}

func (TokenTransfer) TableName() string { return "token_transfers" }
