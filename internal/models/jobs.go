package models

import (
	"time"

	"gorm.io/datatypes"
)

type Job struct {
	ID          uint64         `gorm:"primaryKey;column:id"`
	JobType     string         `gorm:"column:job_type;type:text;not null;index:idx_jobs_type"`
	Payload     datatypes.JSON `gorm:"column:payload;type:jsonb"`
	Status      string         `gorm:"column:status;type:text;not null;default:pending;index:idx_jobs_status"`
	Attempts    int            `gorm:"column:attempts;type:int;default:0"`
	LastError   *string        `gorm:"column:last_error;type:text"`
	ScheduledAt time.Time      `gorm:"column:scheduled_at;type:timestamptz;not null;default:now();index:idx_jobs_scheduled"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (Job) TableName() string {
	return "jobs"
}
