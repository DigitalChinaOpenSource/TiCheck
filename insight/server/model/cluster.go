package model

import (
	"time"
)

type Cluster struct {
	ID            uint `gorm:"primarykey"`
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