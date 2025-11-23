package models

import (
	"time"

	"gorm.io/datatypes"
)

type ContractCall struct {
	ID              int64          `gorm:"primaryKey;column:id"`
	Chain           string         `gorm:"column:chain"`
	TxHash          string         `gorm:"column:tx_hash"`
	ContractAddress *string        `gorm:"column:contract_address"`
	Method          *string        `gorm:"column:method"`
	Args            datatypes.JSON `gorm:"column:args"`
	ReturnData      datatypes.JSON `gorm:"column:return_data"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
}

func (ContractCall) TableName() string {
	return "contract_calls"
}
