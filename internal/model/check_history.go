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

func (ch *CheckHistory) GetHistoryByClusterID(id int, pageSize int, pageNum int, startTime string, endTime string) (map[string]interface{}, error) {
	var chs []CheckHistory
	var total int64
	var err error
	if startTime == "" || endTime == "" {
		err = DbConn.Where("cluster_id = ?", id).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("id desc").Find(&chs).Error
	} else {
		err = DbConn.Where("cluster_id = ? and check_time between ? and ?", id, startTime, endTime).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("id desc").Find(&chs).Error
	}

	if err != nil {
		return nil, err
	}

	DbConn.Model(&chs).Count(&total)

	return map[string]interface{}{
		"page_size": pageSize,
		"page_num":  pageNum,
		"total":     total,
		"data":      chs,
	}, nil
}
