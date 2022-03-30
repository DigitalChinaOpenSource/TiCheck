package model

import (
	"database/sql"
	"time"
)

type Cluster struct {
	ID            uint      `gorm:"primarykey"`
	Name          string    `gorm:"not null"`
	PrometheusURL string    `gorm:"not null"`
	TiDBUsername  string    `gorm:"not null"`
	TiDBPassword  string    `gorm:"not null"`
	LoginPath     string    `gorm:"not null"`
	CreateTime    time.Time `gorm:"not null"`
	Description   string
	Owner         string `gorm:"not null"`
	TiDBVersion   string `gorm:"not null"`
	DashboardURL  string
	GrafanaURL    string
	LastCheckTime time.Time
	ClusterHealth int
}

// RecentWarnings is a struct for QueryRecentWarningsByID
type RecentWarnings struct {
	CheckTime    time.Time `json:"time"`
	WarningItems int       `json:"warnings"`
}

// HistoryWarnings is a struct for QueryHistoryWarningByID
type HistoryWarnings struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

// CheckHistoryInfo is a struct for QueryHistoryInfoByID
type CheckHistoryInfo struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

func (c *Cluster) TableName() string {
	return "clusters"
}

func IsClusterExist(id int) bool {
	var count int64
	err := DbConn.Model(&Cluster{}).Where("id = ?", id).Count(&count).Error
	if err != nil || count < 1 {
		return false
	} else {
		return true
	}
}

// QueryClusterInfoByID query a cluster information by its id
func (c *Cluster) QueryClusterInfoByID(id int) (clusterInfo Cluster, err error) {
	err = DbConn.First(&clusterInfo, id).Error
	if err != nil {
		return clusterInfo, err
	}
	return clusterInfo, nil
}

// QueryClusterList query all clusters under the current user
func (c *Cluster) QueryClusterList(owner string) ([]Cluster, error) {
	var clusterList []Cluster
	err := DbConn.
		Order("create_time asc").
		Where("owner = ?", owner).
		Find(&clusterList).
		Error
	if err != nil {
		return clusterList, err
	}

	return clusterList, nil
}

// CreateCluster create a cluster for current user
func (c *Cluster) CreateCluster() (err error) {
	c.CreateTime = time.Now().Local()
	err = DbConn.Create(&c).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateClusterByID update a cluster by its id
func (c *Cluster) UpdateClusterByID() error {
	updateData := map[string]interface{}{
		"name":           c.Name,
		"prometheus_url": c.PrometheusURL,
		"ti_db_username": c.TiDBUsername,
		"ti_db_password": c.TiDBPassword,
		"ti_db_version":  c.TiDBVersion,
		"description":    c.Description,
		"dashboard_url":  c.DashboardURL,
		"grafana_url":    c.GrafanaURL,
	}
	err := DbConn.Model(&c).
		Where("id = ?", c.ID).
		Updates(updateData).
		Error
	if err != nil {
		return err
	}
	return nil
}

// CheckConn check connectivity to tidb cluster
func (c *Cluster) CheckConn(path string) error {
	DB, err := sql.Open("mysql", path)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	defer DB.Close()

	return nil
}

// QueryHistoryInfoByID query the total number of inspections and
// the total number of inspection items for a cluster by its id.
func (c *Cluster) QueryHistoryInfoByID(id int) (CheckHistoryInfo, error) {
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

// QueryHistoryWarningByID query the warnings items in the cluster in the past period by cluster id
func (c *Cluster) QueryHistoryWarningByID(id int, days int) ([]HistoryWarnings, error) {
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

// QueryRecentWarningsByID query the cluster's recent warnings by its cluster id
func (c *Cluster) QueryRecentWarningsByID(id int) (recentWarnings []RecentWarnings, err error) {
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

// GetLashCheckTime get last check time by cluster id
func (c *Cluster) GetLashCheckTime(id int) (time.Time, error) {
	var lastTime time.Time
	err := DbConn.Model(c).Select("LastCheckTime").Where("id = ?", id).First(&lastTime).Error

	return lastTime, err
}
