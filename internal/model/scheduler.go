package model

import (
	"time"
)

type Scheduler struct {
	ID        	   uint `gorm:"primarykey"`
	Name           string `gorm:"not null"`
	ClusterID      uint   `gorm:"not null"`
	CronExpression string `gorm:"not null"`
	IsEnabled      int    `gorm:"not null"`
	Creator        string `gorm:"not null"`
	RunCount       int    `gorm:"not null"`
	CreateTime     time.Time    `gorm:"not null"`
}

func (s *Scheduler) TableName() string {
	return "schedulers"
}
