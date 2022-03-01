package model

import (
	"time"
)

type User struct {
	ID             uint   		`gorm:"primarykey"`
	UserName       string       `gorm:"unique;not null"`
	USerPassword   string       `gorm:"not null"`
	FullName       string       `gorm:"not null"`
	Email      	   string       `gorm:"not null"`
	IsEnabled      int          `gorm:"not null"` // 0: Disabled, 1: Enabled
	Creator        string       `gorm:"not null"`
	SystemUser     string       `gorm:"not null"` // under which system user is this user created
	CreateTime     time.Time    `gorm:"not null"`
}

func (u *User) VerifyUser() bool {

	return false
}

func (u *User) AddUser() error {
	return nil
}