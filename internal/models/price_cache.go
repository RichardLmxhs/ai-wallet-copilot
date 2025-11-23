package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type PriceCache struct {
	ID           int64           `gorm:"primaryKey;column:id"`
	TokenAddress *string         `gorm:"column:token_address"`
	Chain        string          `gorm:"column:chain"`
	PriceUSD     decimal.Decimal `gorm:"column:price_usd"`
	Source       *string         `gorm:"column:source"`
	TS           time.Time       `gorm:"column:ts"`
}

func (PriceCache) TableName() string {
	return "price_cache"
}
