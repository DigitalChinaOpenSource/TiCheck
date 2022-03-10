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

		dbFile := "store/ticheck.db"

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
		&Probe{},
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
			pm := &ProbeMeta{}
			err = json.Unmarshal(bytes, pm)
			if err != nil || fi.Name() != pm.ID {
				fmt.Println("无法识别的组件包", err)
				continue
			}

			// transfer to db model
			p := Probe{
				ScriptName:  pm.Name,
				FileName:    pm.Main,
				Description: pm.Description,
				IsSystem:    1,
				Creator:     pm.Author.Name,
				CreateTime:  time.Time(pm.CreateTime).Local(),
				UpdateTime:  time.Time(pm.UpdateTime).Local(),
			}
			if len(pm.Tags) > 0 {
				p.Tag = pm.Tags[0]
			}
			if len(pm.Comparators) > 0 {
				p.Comparator = *pm.Comparators[0]
			}
			// insert to db
			DbConn.Create(&p)
		}
	}

}
