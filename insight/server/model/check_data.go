package model

import (
	"time"
)

type CheckData struct {
	ID        	uint          `gorm:"primarykey"`
	CheckTime   time.Time     `gorm:"not null"`
	CheckTag    string        `gorm:"not null"`
	CheckName   string        `gorm:"not null"`
	Comparator  Comparator    `gorm:"embedded"`
	Duration    time.Duration `gorm:"not null"`
	CheckItem   string        `gorm:"not null"`
	CheckValue  string        // null: script no output
	CheckStatus int           `gorm:"not null"` //0:正常, 1:异常_已有, 2:异常_新增
}