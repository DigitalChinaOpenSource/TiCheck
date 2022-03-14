package handler

import (
	"TiCheck/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ClusterHandler struct {
	ClusterInfo      model.Cluster
	CheckHistoryInfo model.CheckHistoryInfo
	RecentWarnings   model.RecentWarnings
	HistoryWarnings  model.HistoryWarnings
	Scheduler        model.Scheduler
}

type QueryHelper struct {
	Url   string `json:"url"`
	Query string `json:"query"`
}

type PrometheusRespMetrics struct {
	Name     string `json:"__name__"`
	Group    string `json:"group"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

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
	ID       int      `json:"id"`
	NodeType string   `json:"type"`
	Instance []string `json:"instance"`
	Count    int      `json:"count"`
}

type ClusterInfoReq struct {
	ID            uint   `json:"id"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	PrometheusUrl string `json:"url"`
	LogUser       string `json:"user"`
	LogPasswd     string `json:"passwd"`
	Description   string `json:"description"`
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
	ID                     uint                    `json:"id"`
	Name                   string                  `json:"name"`
	Version                string                  `json:"version"`
	ClusterOwner           string                  `json:"owner"`
	Description            string                  `json:"description"`
	CreateTime             time.Time               `json:"create_time"`
	LastCheckTime          time.Time               `json:"last_check_time"`
	ClusterHealth          int                     `json:"cluster_health"`
	CheckCount             int                     `json:"check_count"`
	CheckTotal             int                     `json:"check_total"`
	RecentWarningItems     []model.RecentWarnings  `json:"recent_warning_items"`
	WeeklyHistoryWarnings  []model.HistoryWarnings `json:"weekly_history_warnings"`
	YearlyHistoryWarnings  []model.HistoryWarnings `json:"yearly_history_warnings"`
	MonthlyHistoryWarnings []model.HistoryWarnings `json:"monthly_history_warnings"`
}

type ClusterSchedulerListReps struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	Cron       string    `json:"cron"`
	Status     int       `json:"status"`
	Count      int       `json:"count"`
}

func (ch *ClusterHandler) GetClusterList(c *gin.Context) {
	clusterList, err := ch.ClusterInfo.QueryCLusterList()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
		nodeType := []string{"pd", "tidb", "tikv", "tiflash"}
		url := cluster.PrometheusURL + "/api/v1/query"
		nodesInfo, err := getClusterNodesInfo(url, nodeType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		reps.NodesInfo = nodesInfo
		clusterListReps = append(clusterListReps, reps)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   clusterListReps,
	})
}

func (ch *ClusterHandler) GetClusterInfo(c *gin.Context) {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	clusterInfo, err := ch.ClusterInfo.QueryClusterInfoByID(clusterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	checkHistoryInfo, err := ch.CheckHistoryInfo.QueryHistoryInfoByID(clusterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var recentWarnings []model.RecentWarnings
	recentWarnings, err = ch.RecentWarnings.QueryRecentWarningsByID(clusterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	weekly, err := ch.HistoryWarnings.QueryHistoryWarningByID(clusterID, -7)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	monthly, err := ch.HistoryWarnings.QueryHistoryWarningByID(clusterID, -30)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	yearly, err := ch.HistoryWarnings.QueryHistoryWarningByID(clusterID, -365)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var clusterInfoReps = ClusterInfoReps{
		ID:                     clusterInfo.ID,
		Name:                   clusterInfo.Name,
		ClusterOwner:           clusterInfo.Owner,
		Version:                clusterInfo.TiDBVersion,
		ClusterHealth:          clusterInfo.ClusterHealth,
		CreateTime:             clusterInfo.CreateTime,
		Description:            clusterInfo.Description,
		CheckCount:             checkHistoryInfo.Count,
		CheckTotal:             checkHistoryInfo.Total,
		RecentWarningItems:     recentWarnings,
		WeeklyHistoryWarnings:  weekly,
		MonthlyHistoryWarnings: monthly,
		YearlyHistoryWarnings:  yearly,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   clusterInfoReps,
	})
}

func (ch *ClusterHandler) PostClusterInfo(c *gin.Context) {
	clusterInfoReq := &ClusterInfoReq{}
	err := c.BindJSON(clusterInfoReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the request body is wrong",
		})
		return
	}

	err = ch.BuildClusterInfo(clusterInfoReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ch.ClusterInfo.CreateTime = time.Now().Local()

	err = ch.ClusterInfo.CreateCluster()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	return
}

func (ch *ClusterHandler) UpdateClusterInfo(c *gin.Context) {
	clusterInfoReq := &ClusterInfoReq{}
	err := c.BindJSON(clusterInfoReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the request body is wrong",
		})
		return
	}

	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	clusterInfoReq.ID = uint(clusterID)

	err = ch.BuildClusterInfo(clusterInfoReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = ch.ClusterInfo.UpdateClusterByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	return
}

func (ch *ClusterHandler) BuildClusterInfo(req *ClusterInfoReq) error {
	nodeType := []string{"pd", "grafana"}
	url := fmt.Sprintf("http://%s/api/v1/query", req.PrometheusUrl)
	nodes, err := getClusterNodesInfo(url, nodeType)
	if err != nil {
		return err
	}

	var dashboard string
	var grafana string

	if len(nodes[0].Instance) < 1 {
		return errors.New("found no pd server in this tidb cluster")
	}

	dashboard = fmt.Sprintf("http://%s/dashboard", nodes[0].Instance[0])

	pdUrl := fmt.Sprintf("http://%s/pd/api/v1/version", nodes[0].Instance[0])
	version, err := getClusterVersion(pdUrl)
	if err != nil {
		return err
	}

	if len(nodes[1].Instance) > 0 {
		grafana = fmt.Sprintf("http://%s", nodes[1].Instance[0])
	}

	cluster := model.Cluster{
		ID:            req.ID,
		Name:          req.Name,
		PrometheusURL: fmt.Sprintf("http://%s", req.PrometheusUrl),
		TiDBUsername:  req.LogUser,
		TiDBPassword:  req.LogPasswd,
		Description:   req.Description,
		Owner:         req.Owner,
		TiDBVersion:   version,
		GrafanaURL:    grafana,
		DashboardURL:  dashboard,
	}
	ch.ClusterInfo = cluster
	return nil
}

func (ch *ClusterHandler) GetProbeList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "cluster id is invalid",
		})
		return
	}

	var cc model.ClusterChecklist
	cl, err := cc.GetListInfoByClusterID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cl,
	})
	return
}

func (ch *ClusterHandler) GetAddProbeList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "cluster id is invalid",
		})
		return
	}

	var probe model.Probe
	probes, err := probe.GetNotAddedProveListByClusterID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": probes,
	})
	return
}

func (ch *ClusterHandler) AddProbeForCluster(c *gin.Context) {
	cc := &model.ClusterChecklist{}
	err := c.BindJSON(cc)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = cc.AddCheckProbe()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (ch *ClusterHandler) DeleteProbeForCluster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cc := &model.ClusterChecklist{
		ID: uint(id),
	}
	err = cc.DeleteCheckProbe()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (ch *ClusterHandler) GetClusterSchedulerList(c *gin.Context) {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	schedulerList, err := ch.Scheduler.QuerySchedulersByClusterID(clusterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []ClusterSchedulerListReps
	for _, v := range schedulerList {
		item := ClusterSchedulerListReps{
			ID:         v.ID,
			Name:       v.Name,
			CreateTime: v.CreateTime,
			Cron:       v.CronExpression,
			Count:      v.RunCount,
			Status:     v.IsEnabled,
		}
		data = append(data, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}

func getClusterNodesInfo(url string, nodeType []string) (nodesInfo []NodesInfo, err error) {
	queryHelper := QueryHelper{
		Url: url,
	}
	for k, v := range nodeType {
		var instances []string
		queryString := fmt.Sprintf("probe_success{group='%s'}", v)
		queryHelper.Query = queryString
		resp, err := queryHelper.queryWithPrometheus()
		if err != nil {
			return nodesInfo, err
		}
		for _, res := range resp.Data.Result {
			instances = append(instances, res.Metrics.Instance)
		}
		node := NodesInfo{
			ID:       k,
			NodeType: v,
			Instance: instances,
			Count:    len(instances),
		}
		nodesInfo = append(nodesInfo, node)
	}

	return nodesInfo, nil
}

func getClusterVersion(url string) (version string, err error) {
	queryHelper := QueryHelper{
		Url: url,
	}

	jsonMap, err := queryHelper.queryWithPD()
	if err != nil {
		return version, errors.New(fmt.Sprintf("bad request: %s", err))
	}

	version = fmt.Sprintf("%v", jsonMap["version"])
	return version, nil
}

func (h QueryHelper) queryWithPrometheus() (proResp PrometheusResp, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", h.Url, nil)
	if err != nil {
		return proResp, err
	}

	q := req.URL.Query()
	q.Add("query", h.Query)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return proResp, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return proResp, err
	}
	bodyStr := string(body)
	if errJson := json.Unmarshal([]byte(bodyStr), &proResp); errJson != nil {
		return proResp, errJson
	}

	return proResp, nil
}

func (h QueryHelper) queryWithPD() (pdResp map[string]interface{}, err error) {
	resp, err := http.Get(h.Url)
	if err != nil {
		return pdResp, errors.New(fmt.Sprintf("bad request: %s", err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pdResp, errors.New(fmt.Sprintf("parse error: %s", err))
	}

	bodyStr := string(body)
	if errJson := json.Unmarshal([]byte(bodyStr), &pdResp); errJson != nil {
		return pdResp, errors.New(fmt.Sprintf("parse error: %s", err))
	}

	return pdResp, nil
}
