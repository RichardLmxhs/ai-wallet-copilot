package models

import (
	"time"

	"gorm.io/datatypes"
)

type AnalysisResult struct {
	ID            uint64         `gorm:"primaryKey"`
	WalletAddress string         `gorm:"column:wallet_address;type:text;not null;index:idx_analysis_results_wallet"`
	Chain         string         `gorm:"type:text;default:ethereum;not null"`
	RequestID     *string        `gorm:"column:request_id;type:text"`
	ModelName     *string        `gorm:"column:model_name;type:text"`
	PromptHash    *string        `gorm:"column:prompt_hash;type:text"`
	Behavior      *string        `gorm:"column:behavior_summary;type:text"`
	RiskScore     *int           `gorm:"column:risk_score;type:int"`
	RiskDetails   datatypes.JSON `gorm:"column:risk_details;type:jsonb"`
	Suggestions   datatypes.JSON `gorm:"column:suggestions;type:jsonb"`
	RawResponse   datatypes.JSON `gorm:"column:raw_response;type:jsonb"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamptz;default:now();index:idx_analysis_results_created"`
}

func (AnalysisResult) TableName() string { return "analysis_results" }
