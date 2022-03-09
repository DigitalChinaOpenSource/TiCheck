package model

import (
	"crypto/sha1"
	"fmt"
	"time"
)

type User struct {
	ID           uint      `gorm:"primarykey"`
	UserName     string    `gorm:"unique;not null"`
	UserPassword string    `gorm:"not null"`
	FullName     string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	IsEnabled    int       `gorm:"not null"` // 0: Disabled, 1: Enabled
	Creator      string    `gorm:"not null"`
	CreateTime   time.Time `gorm:"not null"`
}

// VerifyUser Check the login permission of the user
func (u *User) VerifyUser() bool {
	var total int64
	err := DbConn.Model(&u).Where("user_name = ? and user_password = ? and is_enabled = 1 ", u.UserName, fmt.Sprintf("%X", sha1.Sum([]byte(u.UserPassword)))).Count(&total).Error
	if err != nil {
		return false
	}

	if total < 1 {
		return false
	}

	return true
}

func (u *User) GetUserInfoByName() error {
	err := DbConn.Select("user_name,full_name,email,is_enabled").Where("user_name = ?", u.UserName).Limit(1).Find(u).Error
	if err != nil {
		return err
	}
	return nil
}
