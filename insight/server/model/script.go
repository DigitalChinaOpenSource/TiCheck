package model

import "time"

type Script struct {
	ID                uint   	 `gorm:"primarykey"`
	ScriptName        string     `gorm:"not null;unique"`
	FileName          string 	 `gorm:"not null;unique"`
	Tag               string  	 `gorm:"not null"`
	Description       string
	DefaultComparator Comparator `gorm:"embedded"`
	IsSystem          int        `gorm:"not null"`
	Creator           string     `gorm:"not null"`
	CreateTime        time.Time  `gorm:"not null"`
	UpdateTime		  time.Time  `gorm:"not null"`
}