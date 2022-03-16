package model

import (
	"time"
)

type CheckData struct {
	ID          uint          `gorm:"primarykey" json:"id"`
	HistoryID   uint          `gorm:"not null" json:"history_id"`
	ClusterID   uint          `gorm:"not null" json:"cluster_id"`
	CheckTime   time.Time     `gorm:"not null" json:"check_time"`
	CheckTag    string        `gorm:"not null" json:"check_tag"`
	CheckName   string        `gorm:"not null" json:"check_name"`
	Comparator    			  `gorm:"embedded" json:",inline"`
	Duration    time.Duration `gorm:"not null" json:"duration"`
	CheckItem   string        `gorm:"not null" json:"check_item"`
	CheckValue  string        `json:"check_value"`                  // null: script no output
	CheckStatus int           `gorm:"not null" json:"check_status"` //0:正常, 1:异常_已有, 2:异常_新增
}

func (cd *CheckData) TableName() string {
	return "check_data"
}

func (cd *CheckData) GetDataByHistoryID(id int) ([]CheckData, error) {
	var cds []CheckData

	err := DbConn.Where("history_id = ?", id).Find(&cds).Error

	if err !=nil {
		return cds, err
	}

	return cds, nil
}