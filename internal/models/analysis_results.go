package models

import (
	"time"

	"gorm.io/datatypes"
)

type AnalysisResult struct {
	ID            uint64         `gorm:"primaryKey"`
	WalletAddress string         `gorm:"type:text;index;not null"`
	Chain         string         `gorm:"type:text;default:ethereum;not null"`
	RequestID     *string        `gorm:"type:text"`
	ModelName     *string        `gorm:"type:text"`
	PromptHash    *string        `gorm:"type:text"`
	Behavior      *string        `gorm:"column:behavior_summary;type:text"`
	RiskScore     *int           `gorm:"type:int"`
	RiskDetails   datatypes.JSON `gorm:"type:jsonb"`
	Suggestions   datatypes.JSON `gorm:"type:jsonb"`
	RawResponse   datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;default:now()"`
}

func (AnalysisResult) TableName() string { return "analysis_results" }
