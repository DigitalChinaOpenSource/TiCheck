package model

import (
	"time"
)

type Cluster struct {
	ID            uint      `gorm:"primarykey"`
	Name          string    `gorm:"not null"`
	PrometheusURL string    `gorm:"not null"`
	TiDBUsername  string    `gorm:"not null"`
	TiDBPassword  string    `gorm:"not null"`
	CreateTime    time.Time `gorm:"not null"`
	Description   string
	Owner         string `gorm:"not null"`
	TiDBVersion   string `gorm:"not null"`
	DashboardURL  string
	GrafanaURL    string
	LastCheckTime time.Time
	ClusterHealth int
}

type RecentWarnings struct {
	CheckTime    time.Time `json:"x"`
	WarningItems int       `json:"y"`
}

type HistoryWarnings struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type CheckHistoryInfo struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

func (c *Cluster) QueryClusterInfoByID(id int) (clusterInfo Cluster, err error) {
	err = DbConn.First(&clusterInfo, id).Error
	if err != nil {
		return clusterInfo, err
	}
	return clusterInfo, nil
}

func (c *Cluster) QueryCLusterList() ([]Cluster, error) {
	var clusterList []Cluster
	err := DbConn.
		Order("create_time asc").
		Find(&clusterList).
		Error
	if err != nil {
		return clusterList, err
	}

	return clusterList, nil
}

func (c *Cluster) CreateCluster() (err error) {
	var cluster Cluster
	err = DbConn.Last(&cluster).Error
	if err != nil {
		return err
	}

	c.ID = cluster.ID + 1

	err = DbConn.Create(&c).Error
	if err != nil {
		return err
	}
	return nil
}

func (chi *CheckHistoryInfo) QueryHistoryInfoByID(id int) (CheckHistoryInfo, error) {
	var checkHistory CheckHistoryInfo
	err := DbConn.Model(&CheckHistory{}).
		Select("count(*) as count,sum(total_items) as total").
		Where("cluster_id", id).
		First(&checkHistory).
		Error
	if err != nil {
		return checkHistory, err
	}
	return checkHistory, nil
}

func (hw *HistoryWarnings) QueryHistoryWarningByID(id int, days int) ([]HistoryWarnings, error) {
	var historyWarnings []HistoryWarnings
	err := DbConn.Model(&CheckData{}).
		Select("check_name || '(' ||check_item || ')' as name,count(*) as total").
		Where("cluster_id = ? AND check_time > ?",
			id,
			time.Now().AddDate(0, 0, days).Format("2006-01-02 15:04:05")).
		Not("check_status", 0).
		Group("name").
		Order("total desc").
		Limit(10).
		Find(&historyWarnings).
		Error
	if err != nil {
		return historyWarnings, err
	}
	return historyWarnings, nil
}

func (rw *RecentWarnings) QueryRecentWarningsByID(id int) (recentWarnings []RecentWarnings, err error) {
	subQuery := DbConn.Model(&CheckHistory{}).
		Select("id,check_time,warning_items").
		Where("cluster_id", id).
		Order("check_time desc").
		Limit(10)
	err = DbConn.Table("(?)", subQuery).
		Order("id").
		Find(&recentWarnings).
		Error
	if err != nil {
		return recentWarnings, err
	}

	return recentWarnings, nil
}
