package executor

import (
	"gorm.io/gorm"
	"time"
)

type Comparator struct {
	// represents way to compare threshold and result
	// 0: NA, 1: eq. 2: g, 3: ge, 4: l, 5: le
	Operator  int
	Threshold string
}

type User struct {
	gorm.Model
	Name       string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	FullName   string `gorm:"not null"`
	Email      string `gorm:"not null"`
	IsEnabled  int    `gorm:"not null"` // 0: Disabled, 1: Enabled
	SystemUser string `gorm:"not null"` // under which system user is this user created
}
type Script struct {
	gorm.Model
	Name              string `gorm:"not null;unique"`
	FileName          string `gorm:"not null;unique"`
	Tag               string `gorm:"not null"`
	Description       string
	DefaultComparator Comparator `gorm:"embedded"`
	IsSystem          int        `gorm:"not null"`
	Creator           string     `gorm:"not null"`
}

type Cluster struct {
	gorm.Model
	Name          string `gorm:"not null;unique"`
	PrometheusURL string `gorm:"not null"`
	TiDBUsername  string `gorm:"not null"`
	TiDBPassword  string `gorm:"not null"`
	Description   string
	Owner         string `gorm:"not null"`
	TiDBVersion   string `gorm:"not null"`
	DashboardURL  string
	GrafanaURL    string
	LastCheckTime time.Time
	ClusterHealth int
}

type ClusterChecklist struct {
	gorm.Model
	ClusterID      uint       `gorm:"not null"`
	ScriptID       uint       `gorm:"not null"`
	IsEnabled      int        `gorm:"not null"`
	Comparator     Comparator `gorm:"embedded"`
	AdditionalArgs string
}

type Scheduler struct {
	gorm.Model
	Name           string `gorm:"not null"`
	ClusterID      uint   `gorm:"not null"`
	CronExpression string `gorm:"not null"`
	IsEnabled      int    `gorm:"not null"`
	Creator        string `gorm:"not null"`
	RunCount       int    `gorm:"not null"`
}

type CheckHistory struct {
	gorm.Model
	CheckTime    time.Time     `gorm:"not null"`
	ClusterID    uint          `gorm:"not null"`
	SchedulerID  uint          // null if run manually
	NormalItems  int           `gorm:"not null"`
	WarningItems int           `gorm:"not null"`
	TotalItems   int           `gorm:"not null"`
	Duration     time.Duration `gorm:"not null"`
}

type CheckData struct {
	gorm.Model
	CheckTime   time.Time     `gorm:"not null"`
	CheckTag    string        `gorm:"not null"`
	CheckName   string        `gorm:"not null"`
	Comparator  Comparator    `gorm:"embedded"`
	Duration    time.Duration `gorm:"not null"`
	CheckItem   string        `gorm:"not null"`
	CheckValue  string        // null: script no output
	CheckStatus int           `gorm:"not null"` //0:正常, 1:异常_已有, 2:异常_新增
}
