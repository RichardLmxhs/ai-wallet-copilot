package models

import (
	"time"

	"gorm.io/datatypes"
)

type Job struct {
	ID          int64          `gorm:"primaryKey;column:id"`
	JobType     string         `gorm:"column:job_type"`
	Payload     datatypes.JSON `gorm:"column:payload"`
	Status      string         `gorm:"column:status"`
	Attempts    int            `gorm:"column:attempts"`
	LastError   *string        `gorm:"column:last_error"`
	ScheduledAt time.Time      `gorm:"column:scheduled_at"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
}

func (Job) TableName() string {
	return "jobs"
}
