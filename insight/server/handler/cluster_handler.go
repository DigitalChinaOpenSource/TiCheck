package handler

import (
	"TiCheck/insight/server/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ClusterHandler struct{}

type CheckHistoryInfo struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

type HistoryWarningItemNum struct {
	CheckTime        time.Time `json:"check_time"`
	TotalWarningItem int       `json:"total_warning_items"`
}

type HistoryWarningItemDetail struct {
	CheckTime       time.Time `json:"check_time"`
	WarningItemName string    `json:"warning_item_name"`
	WarningItemNum  string    `json:"warning_item_num"`
}

type ClusterInfoReps struct {
	ID                uint      `json:"cluster_id"`
	ClusterName       string    `json:"cluster_name"`
	ClusterVersion    string    `json:"cluster_version"`
	ClusterOwner      string    `json:"cluster_owner"`
	Description       string    `json:"cluster_description"`
	CreateTime        time.Time `json:"create_time"`
	LastCheckTime     time.Time `json:"last_check_time"`
	ClusterHealth     int       `json:"cluster_health"`
	CheckCount        int       `json:"check_count"`
	CheckTotal        int       `json:"check_total"`
	TotalWarningItems int       `json:"total_warning_items"`
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
		Select("count(ID) as count,sum(TotalItems) as total").
		Where("ClusterID=", id).
		First(&checkHistory)
	if history.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var historyWarningItemNum HistoryWarningItemNum
	warningItemNum := model.DbConn.Model(&model.CheckData{}).
		Select("CheckTime,count(CheckTime) as TotalWarningItem").
		Where("ClusterID=", id).
		Not("CheckStatus=", 0).
		Group("CheckTime").
		Order("CheckTime desc").
		Limit(10).
		Find(&historyWarningItemNum)
	if warningItemNum.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var historyWarningItemDetail HistoryWarningItemDetail
	itemDetail := model.DbConn.Model(&model.CheckData{}).
		Select("CheckTime,CheckName as WarningItemName,count(CheckName) as WarningItemNum").
		Where("ClusterID=", id).
		Not("CheckStatus=", 0).
		Group("CheckName").
		Order("CheckName desc").
		Find(&historyWarningItemDetail)
	if itemDetail.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var clusterInfoReps = ClusterInfoReps{
		ID:                clusterInfo.ID,
		ClusterName:       clusterInfo.Name,
		ClusterOwner:      clusterInfo.Owner,
		ClusterVersion:    clusterInfo.TiDBVersion,
		ClusterHealth:     clusterInfo.ClusterHealth,
		LastCheckTime:     clusterInfo.LastCheckTime,
		CreateTime:        clusterInfo.CreateTime,
		Description:       clusterInfo.Description,
		CheckCount:        checkHistory.Count,
		CheckTotal:        checkHistory.Total,
		TotalWarningItems: historyWarningItemNum.TotalWarningItem,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   clusterInfoReps,
	})
}
