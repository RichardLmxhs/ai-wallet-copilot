package models

import (
	"time"

	"gorm.io/datatypes"
)

type ContractCall struct {
	ID              uint64         `gorm:"primaryKey;column:id"`
	Chain           string         `gorm:"column:chain;type:text;default:ethereum;not null"`
	TxHash          string         `gorm:"column:tx_hash;type:text;not null;index:idx_contract_calls_tx"`
	ContractAddress *string        `gorm:"column:contract_address;type:text;index:idx_contract_calls_contract"`
	Method          *string        `gorm:"column:method;type:text"`
	Args            datatypes.JSON `gorm:"column:args;type:jsonb"`
	ReturnData      datatypes.JSON `gorm:"column:return_data;type:jsonb"`
	CreatedAt       time.Time      `gorm:"column:created_at;type:timestamptz;default:now()"`
}

func (ContractCall) TableName() string {
	return "contract_calls"
}
