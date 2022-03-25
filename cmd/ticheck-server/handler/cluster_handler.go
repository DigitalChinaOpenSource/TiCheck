package handler

import (
	"TiCheck/cmd/ticheck-server/api"
	"TiCheck/executor"
	"TiCheck/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	ClusterInfo model.Cluster
	Scheduler   model.Scheduler
	Url         string `json:"url"`
	Query       string `json:"query"`
}

// PrometheusRespMetrics is a struct for PrometheusResp metrics
type PrometheusRespMetrics struct {
	Name     string `json:"__name__"`
	Group    string `json:"group"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

// PrometheusRespRes is a struct for PrometheusResp result
type PrometheusRespRes struct {
	Metrics PrometheusRespMetrics `json:"metric"`
	Value   []interface{}         `json:"value"`
}

// PrometheusRespData is a struct for PrometheusRespRes data
type PrometheusRespData struct {
	ResultType string              `json:"resultType"`
	Result     []PrometheusRespRes `json:"result"`
}

// PrometheusResp is a struct for prometheus response
type PrometheusResp struct {
	Status string             `json:"status"`
	Data   PrometheusRespData `json:"data"`
}

// NodesInfo is a struct for tidb cluster's nodes
type NodesInfo struct {
	ID       int      `json:"id"`
	NodeType string   `json:"type"`
	Instance []string `json:"instance"`
	Count    int      `json:"count"`
	Normal   int      `json:"normal"`
}

// ClusterInfoReq is a request for PostClusterInfo
type ClusterInfoReq struct {
	ID            uint     `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	PrometheusUrl string   `json:"url"`
	LogUser       string   `json:"user"`
	LogPasswd     string   `json:"passwd"`
	Description   string   `json:"description"`
	CheckItems    []string `json:"check_items"`
}

// ClusterListResp is a response for GetClusterList
type ClusterListResp struct {
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

// ClusterInfoResp is a response for GetClusterInfo
type ClusterInfoResp struct {
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

// InitialClusterInfoResp is a response for GetInitialClusterInfo
type InitialClusterInfoResp struct {
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	PrometheusUrl string `json:"url"`
	LogUser       string `json:"user"`
	LogPasswd     string `json:"passwd"`
	Description   string `json:"description"`
}

// ClusterSchedulerListResp is a response for GetClusterSchedulerList
type ClusterSchedulerListResp struct {
	Index      int       `json:"index"`
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	Cron       string    `json:"cron"`
	Status     int       `json:"status"`
	Count      int       `json:"count"`
}

// ClusterSchedulerReq is a request for PostClusterScheduler
type ClusterSchedulerReq struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Cron      string `json:"cron"`
	Status    bool   `json:"status"`
	Creator   string `json:"creator"`
	ClusterID string `json:"cluster_id"`
}

// GetClusterList get all clusters of currently log in user
func (ch *ClusterHandler) GetClusterList(c *gin.Context) {
	owner := c.Query("owner")
	clusterList, err := ch.ClusterInfo.QueryClusterList(owner)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var clusterListReps []ClusterListResp
	for _, cluster := range clusterList {
		var reps = ClusterListResp{
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
		nodesInfo, e := GetClusterNodesInfo(url, nodeType)
		if e != nil {
			api.ErrorWithMsg(c, e.Error())
			return
		}
		reps.NodesInfo = nodesInfo
		clusterListReps = append(clusterListReps, reps)
	}

	api.Success(c, "", clusterListReps)
	return
}

// GetClusterInfo get cluster info for cluster info view
func (ch *ClusterHandler) GetClusterInfo(c *gin.Context) {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(clusterID) {
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	clusterInfo, err := ch.ClusterInfo.QueryClusterInfoByID(clusterID)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	checkHistoryInfo, err := ch.ClusterInfo.QueryHistoryInfoByID(clusterID)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var recentWarnings []model.RecentWarnings
	recentWarnings, err = ch.ClusterInfo.QueryRecentWarningsByID(clusterID)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	weekly, err := ch.ClusterInfo.QueryHistoryWarningByID(clusterID, -7)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	monthly, err := ch.ClusterInfo.QueryHistoryWarningByID(clusterID, -30)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	yearly, err := ch.ClusterInfo.QueryHistoryWarningByID(clusterID, -365)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var clusterInfoReps = ClusterInfoResp{
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

	api.Success(c, "", clusterInfoReps)
	return
}

// GetInitialClusterInfo get initial cluster information before cluster updated
func (ch *ClusterHandler) GetInitialClusterInfo(c *gin.Context) {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(clusterID) {
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	clusterInfo, err := ch.ClusterInfo.QueryClusterInfoByID(clusterID)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}
	url := strings.Trim(clusterInfo.PrometheusURL, "http://")
	var initialCluster = InitialClusterInfoResp{
		Name:          clusterInfo.Name,
		Owner:         clusterInfo.Owner,
		PrometheusUrl: url,
		Description:   clusterInfo.Description,
		LogUser:       clusterInfo.TiDBUsername,
		LogPasswd:     clusterInfo.TiDBPassword,
	}

	api.Success(c, "", initialCluster)
	return
}

// PostClusterInfo add a cluster
func (ch *ClusterHandler) PostClusterInfo(c *gin.Context) {
	clusterInfoReq := &ClusterInfoReq{}
	err := c.BindJSON(clusterInfoReq)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	cluster, err := ch.BuildClusterInfo(clusterInfoReq)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cluster.CreateCluster()
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// UpdateClusterInfo update a cluster by cluster id
func (ch *ClusterHandler) UpdateClusterInfo(c *gin.Context) {
	clusterInfoReq := &ClusterInfoReq{}
	err := c.BindJSON(clusterInfoReq)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(clusterID) {
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	clusterInfoReq.ID = uint(clusterID)

	cluster, err := ch.BuildClusterInfo(clusterInfoReq)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cluster.UpdateClusterByID()
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// BuildClusterInfo build a cluster information by ClusterInfoReq for UpdateClusterInfo and PostClusterInfo
func (ch *ClusterHandler) BuildClusterInfo(req *ClusterInfoReq) (cluster model.Cluster, err error) {
	nodeType := []string{"pd", "grafana", "tidb"}
	url := fmt.Sprintf("http://%s/api/v1/query", req.PrometheusUrl)
	nodes, err := GetClusterNodesInfo(url, nodeType)
	if err != nil {
		return cluster, err
	}

	var dashboard string
	var grafana string

	if len(nodes[0].Instance) < 1 {
		return cluster, errors.New("found no pd server in this tidb cluster")
	}

	if len(nodes[2].Instance) < 1 {
		return cluster, errors.New("found no tidb server in this tidb cluster")
	}

	dashboard = fmt.Sprintf("http://%s/dashboard", nodes[0].Instance[0])

	pdUrl := fmt.Sprintf("http://%s/pd/api/v1/version", nodes[0].Instance[0])
	version, err := GetClusterVersion(pdUrl)
	if err != nil {
		return cluster, err
	}

	tidbUrl := fmt.Sprintf("http://%s/info", nodes[2].Instance[0])
	host, portStr, err := GetClusterConnectPath(tidbUrl)
	if err != nil {
		return cluster, err
	}

	if len(nodes[1].Instance) > 0 {
		grafana = fmt.Sprintf("http://%s", nodes[1].Instance[0])
	}

	path := strings.Join([]string{req.LogUser, ":", req.LogPasswd, "@tcp(", host, ":", portStr, ")/information_schema"}, "")
	err = cluster.CheckConn(path)
	if err != nil {
		return cluster, errors.New("tidb database username or password is wrong")
	}

	cluster = model.Cluster{
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
	return cluster, nil
}

func (ch *ClusterHandler) GetProbeList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(id) {
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	var cc model.ClusterChecklist
	cl, err := cc.GetListInfoByClusterID(id)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	api.Success(c, "", cl)
	return
}

func (ch *ClusterHandler) GetAddProbeList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(id) {
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	var probe model.Probe
	probes, err := probe.GetNotAddedProveListByClusterID(id)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.Success(c, "", probes)
	return
}

func (ch *ClusterHandler) AddProbeForCluster(c *gin.Context) {
	cc := &model.ClusterChecklist{}
	err := c.BindJSON(cc)

	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cc.AddCheckProbe()

	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

func (ch *ClusterHandler) ChangeProbeStatus(c *gin.Context) {
	cc := &model.ClusterChecklist{}
	err := c.BindJSON(cc)

	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cc.ChangeProbeStatus()

	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

func (ch *ClusterHandler) UpdateProbeConfig(c *gin.Context) {
	cc := &model.ClusterChecklist{}
	err := c.BindJSON(cc)

	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cc.UpdateProbeConfig()

	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

func (ch *ClusterHandler) DeleteProbeForCluster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	cc := &model.ClusterChecklist{
		ID: uint(id),
	}
	err = cc.DeleteCheckProbe()

	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// GetClusterSchedulerList get cluster scheduler list of this cluster
func (ch *ClusterHandler) GetClusterSchedulerList(c *gin.Context) {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}
	schedulerList, err := ch.Scheduler.QuerySchedulersByClusterID(clusterID)
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var data []ClusterSchedulerListResp
	for k, v := range schedulerList {
		item := ClusterSchedulerListResp{
			Index:      k + 1,
			ID:         v.ID,
			Name:       v.Name,
			CreateTime: v.CreateTime,
			Cron:       v.CronExpression,
			Count:      v.RunCount,
			Status:     v.IsEnabled,
		}
		data = append(data, item)
	}

	api.Success(c, "", data)
	return
}

// PostClusterScheduler add a scheduler for a cluster by cluster id
func (ch *ClusterHandler) PostClusterScheduler(c *gin.Context) {
	schedulerReq := &ClusterSchedulerReq{}
	err := c.BindJSON(schedulerReq)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	isEnabled := 0
	if schedulerReq.Status {
		isEnabled = 1
	}

	clusterID, err := strconv.Atoi(schedulerReq.ClusterID)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	schedulerInfo := model.Scheduler{
		Name:           schedulerReq.Name,
		CronExpression: schedulerReq.Cron,
		Creator:        schedulerReq.Creator,
		ClusterID:      uint(clusterID),
		IsEnabled:      isEnabled,
		CreateTime:     time.Now().Local(),
		RunCount:       0,
	}

	err = schedulerInfo.AddScheduler()
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// UpdateScheduler update a scheduler information by its id
func (ch *ClusterHandler) UpdateScheduler(c *gin.Context) {
	schedulerReq := &ClusterSchedulerReq{}
	err := c.BindJSON(schedulerReq)
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	isEnabled := 0
	if schedulerReq.Status {
		isEnabled = 1
	}

	schedulerInfo := model.Scheduler{
		ID:             uint(schedulerReq.ID),
		Name:           schedulerReq.Name,
		CronExpression: schedulerReq.Cron,
		IsEnabled:      isEnabled,
	}

	err = schedulerInfo.UpdateScheduler()
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// DeleteScheduler delete a scheduler by scheduler id
func (ch *ClusterHandler) DeleteScheduler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	s := model.Scheduler{
		ID: uint(id),
	}
	err = s.DeleteScheduler()
	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// GetClusterNodesInfo get tidb cluster node information through prometheus and node type name
func GetClusterNodesInfo(url string, nodeType []string) (nodesInfo []NodesInfo, err error) {
	queryHelper := ClusterHandler{
		Url: url,
	}
	for k, v := range nodeType {
		var instances []string
		var normal = 0
		queryString := fmt.Sprintf("probe_success{group='%s'}", v)
		queryHelper.Query = queryString
		resp, err := queryHelper.QueryWithPrometheus()
		if err != nil {
			return nodesInfo, err
		}
		for _, res := range resp.Data.Result {
			instances = append(instances, res.Metrics.Instance)
			if res.Value[1] == "1" {
				normal = normal + 1
			}
		}

		resp, err = queryHelper.QueryWithPrometheus()
		if err != nil {
			return nodesInfo, err
		}
		node := NodesInfo{
			ID:       k,
			NodeType: v,
			Instance: instances,
			Count:    len(instances),
			Normal:   normal,
		}
		nodesInfo = append(nodesInfo, node)
	}

	return nodesInfo, nil
}

// GetClusterVersion get tidb cluster version through PD api
func GetClusterVersion(url string) (version string, err error) {
	queryHelper := ClusterHandler{
		Url: url,
	}

	jsonMap, err := queryHelper.QueryWithUrl()
	if err != nil {
		return version, errors.New(fmt.Sprintf("bad request: %s", err))
	}

	version = fmt.Sprintf("%v", jsonMap["version"])
	return version, nil
}

// GetClusterConnectPath get tidb cluster connection path through PD api
func GetClusterConnectPath(url string) (host string, portStr string, err error) {
	queryHelper := ClusterHandler{
		Url: url,
	}

	jsonMap, err := queryHelper.QueryWithUrl()
	if err != nil {
		return host, portStr, errors.New(fmt.Sprintf("bad request: %s", err))
	}

	host = fmt.Sprintf("%v", jsonMap["ip"])
	portStr = fmt.Sprintf("%v", jsonMap["listening_port"])

	return host, portStr, nil
}

// QueryWithPrometheus query information through the interface of prometheus
func (ch *ClusterHandler) QueryWithPrometheus() (proResp PrometheusResp, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ch.Url, nil)
	if err != nil {
		return proResp, err
	}

	q := req.URL.Query()
	q.Add("query", ch.Query)
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

// QueryWithUrl query information through the interface of PD or TiDB
func (ch *ClusterHandler) QueryWithUrl() (result map[string]interface{}, err error) {
	resp, err := http.Get(ch.Url)
	if err != nil {
		return result, errors.New(fmt.Sprintf("bad request: %s", err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, errors.New(fmt.Sprintf("parse error: %s", err))
	}

	bodyStr := string(body)
	if errJson := json.Unmarshal([]byte(bodyStr), &result); errJson != nil {
		return result, errors.New(fmt.Sprintf("parse error: %s", err))
	}

	return result, nil
}

// ExecuteCheck run once check task, response by realtime
func (ch *ClusterHandler) ExecuteCheck(c *gin.Context) {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		api.BadWithMsg(c, "create websocket connection failed!")
		return
	}
	defer ws.Close()

	exe := executor.CreateClusterExecutor(uint(clusterID), 0)

	resultCh := make(chan executor.CheckResult, 10)
	// ctx := context.WithValue(context.Background(), "", "")
	go exe.Execute(resultCh)

	for {
		select {
		case result := <-resultCh:
			ws.WriteJSON(result.Data)
			// fmt.Printf("%+v\n", result)
			if result.IsFinished {
				return
			}
			time.Sleep(time.Second * 1)
		}

	}
}
