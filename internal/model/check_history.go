package model

import (
	"time"
)

type CheckHistory struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CheckTime    time.Time `gorm:"not null" json:"check_time"`
	ClusterID    uint      `gorm:"not null" json:"cluster_id"`
	SchedulerID  uint      `json:"scheduler_id"` // null if run manually
	NormalItems  uint      `gorm:"not null" json:"normal_items"`
	WarningItems uint      `gorm:"not null" json:"warning_items"`
	TotalItems   uint      `gorm:"not null" json:"total_items"`
	Duration     int64     `gorm:"not null" json:"duration"`
	State        string    `gorm:"not null" json:"state"`
}

func (ch *CheckHistory) TableName() string {
	return "check_histories"
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

	DbConn.Model(&chs).Where("cluster_id = ?", id).Count(&total)

	return map[string]interface{}{
		"page_size": pageSize,
		"page_num":  pageNum,
		"total":     total,
		"data":      chs,
	}, nil
}

func (ch *CheckHistory) IsExistRunningByClusterID(id int) (*CheckHistory, error) {
	var his = &CheckHistory{}

	err := DbConn.Where("cluster_id = ? and State = 'running'", id).Order("check_time desc").Limit(1).Find(his).Error

	return his, err
}

func (ch *CheckHistory) UpdateClusterHealthy(id int) error {
	var healthy float64
	err := DbConn.Table(ch.TableName()).Select("(sum(normal_items)*1.0/sum(total_items)*1.0)*100 as healthy").
		Where("cluster_id = ?", id).Limit(1).Find(&healthy).Error
	if err != nil {
		return err
	}
	DbConn.Model(&Cluster{}).Where("id = ?", id).Update("cluster_health", int(healthy))
	return nil
}
