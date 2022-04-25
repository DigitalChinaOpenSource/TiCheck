package handler

import (
	"TiCheck/cmd/ticheck-server/api"
	"TiCheck/executor"
	"TiCheck/internal/model"
	"TiCheck/internal/service"
	"TiCheck/util/logutil"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os/exec"
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
	Normal        bool        `json:"normal"`
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
	LastCheckNormal        uint                    `json:"last_check_normal"`
	LastCheckWarning       uint                    `json:"last_check_warning"`
	ClusterHealth          int                     `json:"cluster_health"`
	HealthUpdateTime       time.Time               `json:"health_update_time"`
	CheckCount             int                     `json:"check_count"`
	CheckTotal             int                     `json:"check_total"`
	TodayCheckCount        int                     `json:"today_check_count"`
	TodayCheckTotal        int                     `json:"today_check_total"`
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
	ClusterID int    `json:"cluster_id"`
}

// GetClusterList get all clusters of currently log in user
func (ch *ClusterHandler) GetClusterList(c *gin.Context) {
	se, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Debug(err.Error())
		api.BadWithMsg(c, err.Error())
		return
	}

	clusterList, err := ch.ClusterInfo.QueryClusterList(se.User.UserName)
	if err != nil {
		logutil.Logger.Error("Can't get cluster list for this user.", zap.String("user", se.User.UserName))
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
			Normal:        true,
		}
		nodeType := []string{"pd", "tidb", "tikv", "tiflash"}
		url := cluster.PrometheusURL + "/api/v1/query"
		nodesInfo, e := GetClusterNodesInfo(url, nodeType)
		if e != nil {
			logutil.Logger.Warn("Can't connect to the prometheus server.",
				zap.String("prometheus url", cluster.PrometheusURL))
			reps.Normal = false
		}
		reps.NodesInfo = nodesInfo
		clusterListReps = append(clusterListReps, reps)
	}

	api.Success(c, "", clusterListReps)
	return
}

// GetClusterInfo get cluster info for cluster info view
func (ch *ClusterHandler) GetClusterInfo(c *gin.Context) {
	_, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Debug(err.Error())
		api.BadWithMsg(c, err.Error())
		return
	}

	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		logutil.Logger.Error("the request body can't be parsed correctly", zap.Error(err))
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(clusterID) {
		logutil.Logger.Debug("cluster does not exist", zap.String("cluster id", strconv.Itoa(clusterID)))
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	clusterInfo, err := ch.ClusterInfo.QueryClusterInfoByID(clusterID)
	if err != nil {
		logutil.Logger.Error("get cluster information error")
		api.ErrorWithMsg(c, err.Error())
		return
	}

	checkHistoryInfo, err := ch.ClusterInfo.QueryHistoryInfoByID(clusterID)
	if err != nil {
		logutil.Logger.Error("get cluster information error")
		api.ErrorWithMsg(c, err.Error())
		return
	}

	todayHistoryInfo, err := ch.ClusterInfo.QueryTodayHistoryInfoByID(clusterID)
	if err != nil {
		logutil.Logger.Error("Error in getting information about today's inspection items of the cluster.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var checkHistory model.CheckHistory
	lastCheck, err := checkHistory.QueryLastHistoryByID(clusterID)
	if err != nil {
		logutil.Logger.Error("Error in getting information about last check of the cluster.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var recentWarnings []model.RecentWarnings
	recentWarnings, err = ch.ClusterInfo.QueryRecentWarningsByID(clusterID)
	if err != nil {
		logutil.Logger.Error("Error in getting information about recent warnings of the cluster.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	weekly, err := ch.ClusterInfo.QueryHistoryWarningByID(clusterID, -7)
	if err != nil {
		logutil.Logger.Error("Error in getting history warnings of the cluster for the past week.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	monthly, err := ch.ClusterInfo.QueryHistoryWarningByID(clusterID, -30)
	if err != nil {
		logutil.Logger.Error("Error in getting history warnings of the cluster for the past month.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	yearly, err := ch.ClusterInfo.QueryHistoryWarningByID(clusterID, -365)
	if err != nil {
		logutil.Logger.Error("Error in getting history warnings of the cluster for the past year.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var clusterInfoReps = ClusterInfoResp{
		ID:                     clusterInfo.ID,
		Name:                   clusterInfo.Name,
		ClusterOwner:           clusterInfo.Owner,
		Version:                clusterInfo.TiDBVersion,
		ClusterHealth:          clusterInfo.ClusterHealth,
		HealthUpdateTime:       clusterInfo.HealthUpdateTime,
		CreateTime:             clusterInfo.CreateTime,
		Description:            clusterInfo.Description,
		LastCheckTime:          clusterInfo.LastCheckTime,
		LastCheckNormal:        lastCheck.NormalItems,
		LastCheckWarning:       lastCheck.WarningItems,
		CheckCount:             checkHistoryInfo.Count,
		CheckTotal:             checkHistoryInfo.Total,
		TodayCheckCount:        todayHistoryInfo.Count,
		TodayCheckTotal:        todayHistoryInfo.Total,
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
	_, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Debug(err.Error())
		api.BadWithMsg(c, err.Error())
		return
	}

	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		logutil.Logger.Error("the request body can't be parsed correctly", zap.Error(err))
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(clusterID) {
		logutil.Logger.Debug("cluster does not exist", zap.String("cluster id", strconv.Itoa(clusterID)))
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	clusterInfo, err := ch.ClusterInfo.QueryClusterInfoByID(clusterID)
	if err != nil {
		logutil.Logger.Error("Error in getting cluster information.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
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
	se, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Debug(err.Error())
		api.BadWithMsg(c, err.Error())
		return
	}

	clusterInfoReq := &ClusterInfoReq{}
	err = c.BindJSON(clusterInfoReq)
	if err != nil {
		logutil.Logger.Error("the request body can't be parsed correctly", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	clusterInfoReq.Owner = se.User.UserName

	cluster, err := ch.BuildClusterInfo(clusterInfoReq)
	if err != nil {
		logutil.Logger.Error("Failed to get cluster details from prometheus.",
			zap.String("prometheus url", clusterInfoReq.PrometheusUrl))
		api.BadWithMsg(c, err.Error())
		return
	}

	checkList, err := ch.InitialCheckList(clusterInfoReq)
	if err != nil {
		logutil.Logger.Error("Failed to initialize check list.", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	if err = cluster.CreateClusterInTx(checkList); err != nil {
		logutil.Logger.Error("Failed to add cluster.", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// UpdateClusterInfo update a cluster by cluster id
func (ch *ClusterHandler) UpdateClusterInfo(c *gin.Context) {
	se, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Debug(err.Error())
		api.BadWithMsg(c, err.Error())
		return
	}

	clusterInfoReq := &ClusterInfoReq{}
	err = c.BindJSON(clusterInfoReq)
	if err != nil {
		logutil.Logger.Error("the request body can't be parsed correctly", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		logutil.Logger.Error("the cluster id can't be parsed correctly", zap.Error(err))
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(clusterID) {
		logutil.Logger.Debug("cluster does not exist",
			zap.String("cluster id", strconv.Itoa(clusterID)))
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	clusterInfoReq.ID = uint(clusterID)
	clusterInfoReq.Owner = se.User.UserName

	cluster, err := ch.BuildClusterInfo(clusterInfoReq)
	if err != nil {
		logutil.Logger.Error("Failed to get cluster details from prometheus.",
			zap.String("prometheus url", clusterInfoReq.PrometheusUrl))
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cluster.UpdateClusterByID()
	if err != nil {
		logutil.Logger.Error("Failed to update  the cluster.",
			zap.String("cluster id", strconv.Itoa(clusterID)))
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

	lastID, err := cluster.QueryLastID()
	if err != nil {
		return cluster, errors.New("get last cluster id error")
	}
	req.ID = lastID + 1
	loginPath := fmt.Sprintf("ticheck_%d", req.ID)

	if _, err = exec.Command("sh", "../../logpath.sh", loginPath, req.LogUser, host, portStr, req.LogPasswd).Output(); err != nil {
		return cluster, errors.New("failed to add login-path for mysql")
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
		LoginPath:     loginPath,
	}
	return cluster, nil
}

// InitialCheckList initialize the cluster checklist when adding a cluster
func (ch *ClusterHandler) InitialCheckList(req *ClusterInfoReq) (checkList []model.ClusterChecklist, err error) {

	for _, v := range req.CheckItems {
		p := model.Probe{
			ID: v,
		}
		if err = p.GetByID(); err != nil && p.IsSystem != 0 {
			return checkList, err
		}

		initialCom := model.Comparator{
			Operator:  p.Operator,
			Threshold: p.Threshold,
			Arg:       p.Arg,
		}

		checkItem := model.ClusterChecklist{
			ClusterID:  req.ID,
			ProbeID:    v,
			IsEnabled:  0,
			Comparator: initialCom,
		}
		checkList = append(checkList, checkItem)
	}
	return checkList, nil
}

func (ch *ClusterHandler) GetProbeList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logutil.Logger.Error("the request body can't be parsed correctly", zap.Error(err))
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(id) {
		logutil.Logger.Debug("cluster does not exist",
			zap.String("cluster id", strconv.Itoa(id)))
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	var cc model.ClusterChecklist
	cl, err := cc.GetListInfoByClusterID(id)
	if err != nil {
		logutil.Logger.Warn("Failed to get this cluster check list information.",
			zap.String("cluster id", strconv.Itoa(id)))
		api.BadWithMsg(c, err.Error())
		return
	}

	api.Success(c, "", cl)
	return
}

func (ch *ClusterHandler) GetAddProbeList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, "cluster id is invalid")
		return
	}

	if !model.IsClusterExist(id) {
		logutil.Logger.Warn("this cluster does not exist",
			zap.String("cluster id", strconv.Itoa(id)))
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	var probe model.Probe
	probes, err := probe.GetNotAddedProveListByClusterID(id)
	if err != nil {
		logutil.Logger.Warn("failed to get the probe list that has not been added to the cluster",
			zap.String("cluster id", strconv.Itoa(id)))
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
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cc.AddCheckProbe()

	if err != nil {
		logutil.Logger.Warn("failed to add check probe to the cluster",
			zap.String("cluster id", strconv.Itoa(int(cc.ClusterID))))
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
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cc.ChangeProbeStatus()

	if err != nil {
		logutil.Logger.Warn("failed to change the probe status",
			zap.String("probe name", cc.ProbeID))
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
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	err = cc.UpdateProbeConfig()

	if err != nil {
		logutil.Logger.Warn("failed to update probe config",
			zap.String("probe name", cc.ProbeID))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

func (ch *ClusterHandler) DeleteProbeForCluster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	cc := &model.ClusterChecklist{
		ID: uint(id),
	}
	err = cc.DeleteCheckProbe()

	if err != nil {
		logutil.Logger.Warn("failed to delete the cluster probe item",
			zap.String("cluster id", strconv.Itoa(id)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	api.S(c)
	return
}

// GetClusterSchedulerList get cluster scheduler list of this cluster
func (ch *ClusterHandler) GetClusterSchedulerList(c *gin.Context) {
	_, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Debug("the token is invalid.", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}
	schedulerList, err := ch.Scheduler.QuerySchedulersByClusterID(clusterID)
	if err != nil {
		logutil.Logger.Warn("failed to get the cluster scheduler list",
			zap.String("cluster id", strconv.Itoa(clusterID)))
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
	se, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Warn("the token is invalid.", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	schedulerReq := &ClusterSchedulerReq{}
	err = c.BindJSON(schedulerReq)
	if err != nil {
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	isEnabled := 0
	if schedulerReq.Status {
		isEnabled = 1
	}

	schedulerInfo := model.Scheduler{
		Name:           schedulerReq.Name,
		CronExpression: schedulerReq.Cron,
		Creator:        se.User.UserName,
		ClusterID:      uint(schedulerReq.ClusterID),
		IsEnabled:      isEnabled,
		CreateTime:     time.Now().Local(),
		RunCount:       0,
	}

	err = schedulerInfo.AddScheduler()
	if err != nil {
		logutil.Logger.Warn("failed to add scheduler to the cluster",
			zap.String("cluster id", strconv.Itoa(schedulerReq.ClusterID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	// After successfully adding scheduler,according to its IsEnabled value,
	// to decide whether to add cron job for it
	if schedulerInfo.IsEnabled == 1 {
		err = service.CronService.AddTask(schedulerInfo)
		if err != nil {
			logutil.Logger.Warn("failed to add scheduler task to the cron service",
				zap.String("scheduler id", strconv.Itoa(int(schedulerInfo.ID))))
			api.ErrorWithMsg(c, "Failed to run scheduled task")
			return
		}
	}

	api.S(c)
	return
}

// UpdateScheduler update a scheduler information by its id
func (ch *ClusterHandler) UpdateScheduler(c *gin.Context) {
	_, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Warn("the token is invalid.", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	schedulerReq := &ClusterSchedulerReq{}
	err = c.BindJSON(schedulerReq)
	if err != nil {
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
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
		logutil.Logger.Warn("failed to update scheduler information",
			zap.String("scheduler id", strconv.Itoa(schedulerReq.ID)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	// todo there is panic err when update scheduler
	//for _, v := range service.CronService.Tasks {
	//	if v.SchedulerID == schedulerInfo.ID {
	//		service.CronService.RemoveTask(v)
	//	}
	//}

	//if schedulerReq.Status {
	//	err = service.CronService.AddTask(schedulerInfo)
	//	if err != nil {
	//		api.ErrorWithMsg(c, "Failed to update scheduled task")
	//		return
	//	}
	//}

	api.S(c)
	return
}

// DeleteScheduler delete a scheduler by scheduler id
func (ch *ClusterHandler) DeleteScheduler(c *gin.Context) {
	_, err := ch.AccessToken(c)
	if err != nil {
		logutil.Logger.Warn("the token is invalid.", zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logutil.Logger.Warn("the request body can't be parsed correctly",
			zap.Error(err))
		api.BadWithMsg(c, err.Error())
		return
	}

	s := model.Scheduler{
		ID: uint(id),
	}
	err = s.DeleteScheduler()
	if err != nil {
		logutil.Logger.Error("failed to delete the scheduler",
			zap.String("scheduler id", strconv.Itoa(id)))
		api.ErrorWithMsg(c, err.Error())
		return
	}

	for _, v := range service.CronService.Tasks {
		if v.SchedulerID == s.ID {
			service.CronService.RemoveTask(v)
		}
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
			ws.WriteJSON(result)
			// fmt.Printf("%+v\n", result)
			if result.IsFinished {
				return
			}
			time.Sleep(time.Second * 1)
		}

	}
}

// GetExecuteInfo Get some information before execute check.
// include total check times, last check time and check item count for each tag
func (ch *ClusterHandler) GetExecuteInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		api.BadWithMsg(c, err.Error())
		return
	}

	if !model.IsClusterExist(id) {
		api.BadWithMsg(c, "cluster does not exist")
		return
	}

	var cluster model.Cluster
	lastCheckTime, err := cluster.GetLashCheckTime(id)

	historyInfo, err := cluster.QueryHistoryInfoByID(id)

	if err != nil {
		api.ErrorWithMsg(c, err.Error())
		return
	}

	var ccl model.ClusterChecklist
	itemInfo := ccl.GetEnabledCheckListTagGroup(id)

	api.Success(c, "", map[string]interface{}{
		"check_times":     historyInfo.Count,
		"last_check_time": lastCheckTime,
		"check_items":     itemInfo,
	})
	return
}

func (ch *ClusterHandler) AccessToken(c *gin.Context) (se *Session, err error) {
	token, ok := c.Request.Header["Access-Token"]
	if !ok || len(token) < 1 {
		return se, errors.New("the token is invalid")
	}
	se = SessionHelper.getSessionByToken(token[0])
	if se == nil {
		return se, errors.New("can't get session user by this token")
	}
	return se, nil
}
