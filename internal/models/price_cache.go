package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type PriceCache struct {
	ID           uint64                  `gorm:"primaryKey;column:id"`
	TokenAddress *string                 `gorm:"column:token_address;type:text"`
	Chain        string                  `gorm:"column:chain;type:text;default:ethereum;not null"`
	PriceUSD     *decimal.Decimal        `gorm:"column:price_usd;type:numeric(30,8)"`
	Source       *string                 `gorm:"column:source;type:text"`
	Symbol       *string                 `gorm:"column:symbol;type:text"`
	Logo         *string                 `gorm:"column:logo;type:text"`
	UpdatedAt    time.Time               `gorm:"column:updated_at;type:timestamptz;default:now()"`
	Metadata     *map[string]interface{} `gorm:"column:metadata;type:jsonb"`
}

func (PriceCache) TableName() string {
	return "price_cache"
}

// 添加索引配置
func (PriceCache) Indexes() [][]string {
	return [][]string{
		{"token_address", "updated_at DESC"}, // idx_price_token_time
	}
}
