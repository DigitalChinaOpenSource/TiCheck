package model

import (
	"time"
)

type CheckHistory struct {
	ID           uint          `gorm:"primarykey"`
	CheckTime    time.Time     `gorm:"not null"`
	ClusterID    uint          `gorm:"not null"`
	SchedulerID  uint          // null if run manually
	NormalItems  int           `gorm:"not null"`
	WarningItems int           `gorm:"not null"`
	TotalItems   int           `gorm:"not null"`
	Duration     time.Duration `gorm:"not null"`
}