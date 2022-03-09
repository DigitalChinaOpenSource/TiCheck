package handler

import (
	"TiCheck/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ClusterHandler struct{}

type CheckHistoryInfo struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

type RecentWarningItems struct {
	CheckTime    time.Time `json:"check_time"`
	WarningItems int       `json:"total_warning_items"`
}

type HistoryWarnings struct {
	CheckTime      time.Time `json:"check_time"`
	CheckName      string    `json:"check_name"`
	CheckItem      string    `json:"check_item"`
	WarningItemNum string    `json:"warning_item_num"`
}

type ClusterListReps struct {
	ID            uint      `json:"id"`
	Name          string    `json:"cluster_name"`
	Description   string    `json:"description"`
	DashboardUrl  string    `json:"dashboard_url"`
	GrafanaUrl    string    `json:"grafana_url"`
	CreateTime    time.Time `json:"create_time"`
	LastCheckTime time.Time `json:"last_check_time"`
}

type ClusterInfoReps struct {
	ID                     uint                 `json:"id"`
	Name                   string               `json:"name"`
	Version                string               `json:"version"`
	ClusterOwner           string               `json:"owner"`
	Description            string               `json:"description"`
	CreateTime             time.Time            `json:"create_time"`
	LastCheckTime          time.Time            `json:"last_check_time"`
	ClusterHealth          int                  `json:"cluster_health"`
	CheckCount             int                  `json:"check_count"`
	CheckTotal             int                  `json:"check_total"`
	RecentWarningItems     []RecentWarningItems `json:"recent_warning_items"`
	WeeklyHistoryWarnings  []HistoryWarnings    `json:"weekly_history_warnings"`
	YearlyHistoryWarnings  []HistoryWarnings    `json:"yearly_history_warnings"`
	MonthlyHistoryWarnings []HistoryWarnings    `json:"monthly_history_warnings"`
}

func (ch *ClusterHandler) GetClusterList(c *gin.Context) {
	var clusterList []model.Cluster
	res := model.DbConn.
		Order("create_time asc").
		Find(&clusterList)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var clusterListReps []ClusterListReps
	for _, cluster := range clusterList {
		var reps = ClusterListReps{
			ID:            cluster.ID,
			Name:          cluster.Name,
			Description:   cluster.Description,
			DashboardUrl:  cluster.DashboardURL,
			GrafanaUrl:    cluster.GrafanaURL,
			CreateTime:    cluster.CreateTime,
			LastCheckTime: cluster.LastCheckTime,
		}
		clusterListReps = append(clusterListReps, reps)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   clusterListReps,
	})
}

func (ch *ClusterHandler) GetClusterInfo(c *gin.Context) {
	id := c.Param("id")
	var clusterInfo model.Cluster
	res := model.DbConn.First(&clusterInfo, id)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var checkHistory CheckHistoryInfo
	history := model.DbConn.Model(&model.CheckHistory{}).
		Select("count(*) as count,sum(total_items) as total").
		Where("cluster_id", id).
		First(&checkHistory)
	if history.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var recentWarningItems []RecentWarningItems
	warningItemNum := model.DbConn.Model(&model.CheckHistory{}).
		Select("check_time,warning_items").
		Where("cluster_id", id).
		Order("check_time asc").
		Limit(10).
		Find(&recentWarningItems)
	if warningItemNum.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var weekly []HistoryWarnings
	weeklyItems := model.DbConn.Model(&model.CheckData{}).
		Select("check_time,check_name,check_item,count(*) as WarningItemNum").
		Where("cluster_id = ? AND check_time > ?", id, time.Now().AddDate(0, 0, -7)).
		Not("check_status", 0).
		Group("check_name,check_item").
		Order("WarningItemNum desc").
		Limit(10).
		Find(&weekly)
	if weeklyItems.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var monthly []HistoryWarnings
	monthlyItems := model.DbConn.Model(&model.CheckData{}).
		Select("check_time,check_name,check_item,count(*) as WarningItemNum").
		Where("cluster_id = ? AND check_time > ?", id, time.Now().AddDate(0, 1, 0)).
		Not("check_status", 0).
		Group("check_name,check_item").
		Order("WarningItemNum desc").
		Limit(10).
		Find(&monthly)
	if monthlyItems.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var yearly []HistoryWarnings
	yearlyItems := model.DbConn.Model(&model.CheckData{}).
		Select("check_time,check_name,check_item,count(*) as WarningItemNum").
		Where("cluster_id = ? AND check_time > ?", id, time.Now().AddDate(1, 0, 0)).
		Not("check_status", 0).
		Group("check_name,check_item").
		Order("WarningItemNum desc").
		Limit(10).
		Find(&yearly)
	if yearlyItems.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var clusterInfoReps = ClusterInfoReps{
		ID:                     clusterInfo.ID,
		Name:                   clusterInfo.Name,
		ClusterOwner:           clusterInfo.Owner,
		Version:                clusterInfo.TiDBVersion,
		ClusterHealth:          clusterInfo.ClusterHealth,
		LastCheckTime:          clusterInfo.LastCheckTime,
		CreateTime:             clusterInfo.CreateTime,
		Description:            clusterInfo.Description,
		CheckCount:             checkHistory.Count,
		CheckTotal:             checkHistory.Total,
		RecentWarningItems:     recentWarningItems,
		WeeklyHistoryWarnings:  weekly,
		MonthlyHistoryWarnings: monthly,
		YearlyHistoryWarnings:  yearly,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   clusterInfoReps,
	})
}
