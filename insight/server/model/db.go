package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DbConn *gorm.DB

func InitDB() error {
	if DbConn == nil {
		db, err := gorm.Open(sqlite.Open("report/report.db"), &gorm.Config{})
		if err != nil {
			return err
		}

		DbConn = db
	}

	// 表结构未变更时，建议注释此行
	err := DbConn.AutoMigrate(&User{}, &Script{}, &Scheduler{}, &Cluster{}, &CheckHistory{}, &CheckData{}, &ClusterChecklist{})

	if err != nil {
		return err
	}

	return nil
}
