package handler

import (
	"TiCheck/insight/server/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ClusterHandler struct{}

type CheckHistoryInfo struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

type RecentWarningItems struct {
	CheckTime    time.Time `json:"x"`
	WarningItems int       `json:"y"`
}

type HistoryWarnings struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type PrometheusRespMetrics struct{}

type PrometheusRespRes struct {
	Metrics PrometheusRespMetrics `json:"metric"`
	Value   []interface{}         `json:"value"`
}

type PrometheusRespData struct {
	ResultType string              `json:"resultType"`
	Result     []PrometheusRespRes `json:"result"`
}

type PrometheusResp struct {
	Status string             `json:"status"`
	Data   PrometheusRespData `json:"data"`
}

type NodesInfo struct {
	ID       int    `json:"id"`
	NodeType string `json:"type"`
	Count    int    `json:"count"`
}

type ClusterListReps struct {
	ID            uint        `json:"id"`
	Name          string      `json:"cluster_name"`
	Description   string      `json:"description"`
	DashboardUrl  string      `json:"dashboard_url"`
	GrafanaUrl    string      `json:"grafana_url"`
	PrometheusUrl string      `json:"prometheus_url"`
	NodesInfo     []NodesInfo `json:"nodes"`
	CreateTime    time.Time   `json:"create_time"`
	LastCheckTime time.Time   `json:"last_check_time"`
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
			PrometheusUrl: cluster.PrometheusURL,
			DashboardUrl:  cluster.DashboardURL,
			GrafanaUrl:    cluster.GrafanaURL,
			CreateTime:    cluster.CreateTime,
			LastCheckTime: cluster.LastCheckTime,
		}
		reps.NodesInfo, _ = getCLusterNodes(cluster.PrometheusURL)
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
	subQuery := model.DbConn.Model(&model.CheckHistory{}).
		Select("id,check_time,warning_items").
		Where("cluster_id", id).
		Order("check_time desc").
		Limit(10)
	warningItemNum := model.DbConn.Table("(?)", subQuery).
		Order("id").
		Find(&recentWarningItems)
	if warningItemNum.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var weekly []HistoryWarnings
	weeklyItems := model.DbConn.Model(&model.CheckData{}).
		Select("check_name || '(' ||check_item || ')' as name,count(*) as total").
		Where("cluster_id = ? AND check_time > ?", id, time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")).
		Not("check_status", 0).
		Group("name").
		Order("total desc").
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
		Select("check_name || '(' ||check_item || ')' as name,count(*) as total").
		Where("cluster_id = ? AND check_time > ?", id, time.Now().AddDate(0, -1, 0).Format("2006-01-02 15:04:05")).
		Not("check_status", 0).
		Group("name").
		Order("total desc").
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
		Select("check_name || '(' ||check_item || ')' as name,count(*) as total").
		Where("cluster_id = ? AND check_time > ?", id, time.Now().AddDate(-1, 0, 0).Format("2006-01-02 15:04:05")).
		Not("check_status", 0).
		Group("name").
		Order("total desc").
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

func getCLusterNodes(url string) (nodes []NodesInfo, err error) {
	nodeType := []string{"pd", "tidb", "tikv", "tiflash"}
	for k, v := range nodeType {
		num, err := queryNodeNum(url, v)
		if err != nil {
			return nodes, err
		}
		var item = NodesInfo{
			ID:       k,
			NodeType: v,
			Count:    num,
		}
		nodes = append(nodes, item)
	}
	return nodes, nil
}

func queryNodeNum(url string, nodeType string) (num int, err error) {
	var nodeNum int
	queryString := fmt.Sprintf("count(probe_success{group='%s'})", nodeType)
	client := &http.Client{}
	url = url + "api/v1/query"
	req, _ := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("query", queryString)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nodeNum, errors.New(fmt.Sprintf("bad request: %s", err))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nodeNum, errors.New(fmt.Sprintf("parse error: %s", err))
	}
	bodyStr := string(body)
	var jsonMap PrometheusResp
	if errJson := json.Unmarshal([]byte(bodyStr), &jsonMap); errJson != nil {
		return nodeNum, errors.New(fmt.Sprintf("parse error: %s", err))
	}
	if len(jsonMap.Data.Result) < 1 {
		return nodeNum, nil
	}
	for _, v := range jsonMap.Data.Result[0].Value {
		switch value := v.(type) {
		case string:
			nodeNum, err = strconv.Atoi(value)
			if err != nil {
				return nodeNum, errors.New(fmt.Sprintf("parse error: %s", err))
			}
		}
	}
	return nodeNum, nil
}
