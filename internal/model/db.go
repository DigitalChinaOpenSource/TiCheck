package model

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DbConn *gorm.DB

func InitDB() error {
	newDB := false
	if DbConn == nil {

		dbFile := "../../store/ticheck.db"

		_, err := os.Stat(dbFile)
		if os.IsNotExist(err) {
			newDB = true
			file, err := os.Create(dbFile)
			if err != nil {
				return err
			}
			file.Close()
		}

		db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			return err
		}

		DbConn = db
	}

	// 表结构未变更时，建议注释此行
	err := DbConn.AutoMigrate(
		&User{},
		&Script{},
		&Scheduler{},
		&Cluster{},
		&CheckHistory{},
		&CheckData{},
		&ClusterChecklist{},
	)

	if err != nil {
		return err
	}

	if newDB {
		SetupSeedData()
	}

	return nil
}

// init seed data, such as admin user/local script...
func SetupSeedData() {

	// admin user
	admin := User{
		UserName:     "admin",
		UserPassword: fmt.Sprintf("%X", sha1.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte("admin")))))),
		FullName:     "admin",
		Email:        "admin@ticheck.com",
		IsEnabled:    1,
		Creator:      "system",
		CreateTime:   time.Now(),
	}
	result := DbConn.Create(&admin)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	// local script
	scriptDir := "../../probes/local/"
	fileInfos, err := ioutil.ReadDir(scriptDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fileInfos {
		if fi.IsDir() {
			//read package.json
			packageJSON := scriptDir + fi.Name() + "/package.json"
			bytes, err := ioutil.ReadFile(packageJSON)
			if err != nil {
				fmt.Println("读取组件信息失败", err)
				continue
			}
			sm := &ScriptMeta{}
			err = json.Unmarshal(bytes, sm)
			if err != nil || fi.Name() != sm.ID {
				fmt.Println("无法识别的组件包", err)
				continue
			}

			// transfer to db model
			s := Script{
				ScriptName:  sm.Name,
				FileName:    sm.Main,
				Description: sm.Description,
				IsSystem:    1,
				Creator:     sm.Author.Name,
				CreateTime:  sm.CreateTime,
				UpdateTime:  sm.UpdateTime,
			}
			if len(sm.Tags) > 0 {
				s.Tag = sm.Tags[0]
			}
			if len(sm.Comparators) > 0 {
				s.DefaultComparator = *sm.Comparators[0]
			}
			// insert to db
			DbConn.Create(&s)
		}
	}

}
