package models

import (
	"time"

	"gorm.io/datatypes"
)

type WalletNFT struct {
	ID              uint64         `gorm:"primaryKey"`
	WalletAddress   string         `gorm:"column:wallet_address;type:text;not null;index:idx_wallet_nfts_wallet"`
	Chain           string         `gorm:"type:text;default:ethereum;not null"`
	ContractAddress string         `gorm:"column:contract_address;type:text;not null;index:idx_wallet_nfts_contract"`
	TokenID         string         `gorm:"column:token_id;type:text;not null;index:idx_wallet_nfts_token_id"`
	TokenType       string         `gorm:"column:token_type;type:text;default:ERC721"`
	Metadata        datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;default:now()"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;default:now();index:idx_wallet_nfts_updated"`
}

func (WalletNFT) TableName() string { return "wallet_nfts" }
