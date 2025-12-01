package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type WalletToken struct {
	ID            uint64          `gorm:"primaryKey"`
	WalletAddress string          `gorm:"column:wallet_address;type:text;not null;index:idx_wallet_tokens_wallet,unique:ux_wallet_tokens_wallet_token"`
	Chain         string          `gorm:"type:text;default:ethereum;not null"`
	TokenAddress  string          `gorm:"column:token_address;type:text;not null;index:idx_wallet_tokens_token,unique:ux_wallet_tokens_wallet_token"`
	TokenType     string          `gorm:"column:token_type;type:text;default:ERC20"`
	Balance       decimal.Decimal `gorm:"type:numeric(50,18);default:0"`
	Decimals      int             `gorm:"type:int;default:18"`
	Symbol        *string         `gorm:"type:text"`
	Name          *string         `gorm:"type:text"`
	Metadata      datatypes.JSON  `gorm:"type:jsonb"`
	CreatedAt     time.Time       `gorm:"type:timestamptz;default:now()"`
	UpdatedAt     time.Time       `gorm:"type:timestamptz;default:now();index:idx_wallet_tokens_updated"`
}

func (WalletToken) TableName() string { return "wallet_tokens" }
