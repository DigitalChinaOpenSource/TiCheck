package model

import (
	"time"
)

type CheckHistory struct {
	ID           uint          `gorm:"primarykey" json:"id,omitempty"`
	CheckTime    time.Time     `gorm:"not null" json:"check_time"`
	ClusterID    uint          `gorm:"not null" json:"cluster_id,omitempty"`
	SchedulerID  uint          `json:"scheduler_id,omitempty"` // null if run manually
	NormalItems  int           `gorm:"not null" json:"normal_items,omitempty"`
	WarningItems int           `gorm:"not null" json:"warning_items,omitempty"`
	TotalItems   int           `gorm:"not null" json:"total_items,omitempty"`
	Duration     time.Duration `gorm:"not null" json:"duration,omitempty"`
}

func (ch *CheckHistory) GetHistoryByClusterID(id int) ([]CheckHistory,error) {
	var chs []CheckHistory
	err := DbConn.Find(&chs).Where("cluster_id = ?", id).Error

	if err != nil {
		return chs, err
	}

	DbConn.Model(&chs)

	return chs, nil
}