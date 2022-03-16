package model

import (
	"time"
)

type Scheduler struct {
	ID             uint      `gorm:"primarykey"`
	Name           string    `gorm:"not null"`
	ClusterID      uint      `gorm:"not null"`
	CronExpression string    `gorm:"not null"`
	IsEnabled      int       `gorm:"not null"`
	Creator        string    `gorm:"not null"`
	RunCount       int       `gorm:"not null"`
	CreateTime     time.Time `gorm:"not null"`
}

func (s *Scheduler) QuerySchedulersByClusterID(id int) (schedulerList []Scheduler, err error) {
	err = DbConn.
		Where("cluster_id = ?", id).
		Order("create_time asc").
		Find(&schedulerList).
		Error
	if err != nil {
		return schedulerList, err
	}

	return schedulerList, nil
}

func (s *Scheduler) AddScheduler() error {
	err := DbConn.Create(&s).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) TableName() string {
	return "schedulers"
}
