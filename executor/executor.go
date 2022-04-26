package executor

import (
	"TiCheck/config"
	"TiCheck/internal/model"
	"TiCheck/util/logutil"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type CheckResult struct {
	CheckID    uint              `json:"check_id"`
	IsFinished bool              `json:"is_finished"`
	IsTimeout  bool              `json:"is_timeout"`
	IsConfict  bool              `json:"is_confict"`
	Err        string            `json:"err"` // script level error
	Data       []model.CheckData `json:"data"`
}

type Executor interface {
	Execute(result chan CheckResult) // executes one round of check
}

type ClusterExecutor struct {
	ClusterID   uint
	SchedulerID uint
	Prometheus  string
	LoginPath   string
	HistoryID   uint
	CheckList   []model.CheckListInfo
}

type ExecutorContext struct {
	cluster   *ClusterExecutor
	checkInfo model.CheckListInfo

	wg      *sync.WaitGroup
	counter *ExecutorCounter
}

type ExecutorCounter struct {
	mutex sync.Mutex

	normalItems  uint
	warningItems uint
}

// Execute the cluster exexute once check, rc is check result
func (ce *ClusterExecutor) Execute(rc chan CheckResult) {
	logutil.Logger.Info("cluster check start.", zap.String("cluster id", strconv.Itoa(int(ce.ClusterID))))
	// check if historiy in running
	history, err := (&model.CheckHistory{}).IsExistRunningByClusterID(int(ce.ClusterID))
	if err == nil && history.ID > 0 {
		msg := fmt.Sprintf("Execute confict, There is a running check task with a start time of %s", history.CheckTime)
		rc <- CheckResult{
			IsConfict:  true,
			IsFinished: true,
			Err:        msg,
		}
		logutil.Logger.Error(msg, zap.String("cluster id", strconv.Itoa(int(ce.ClusterID))))
		return
	}

	begin := time.Now()
	// create check history first
	his := model.CheckHistory{
		CheckTime:   begin,
		ClusterID:   ce.ClusterID,
		SchedulerID: ce.SchedulerID,
		State:       "running",
	}
	if err := model.DbConn.Create(&his).Error; err != nil {
		result := CheckResult{IsFinished: true}
		result.Err = fmt.Sprintf("create check history error: %s", err.Error())
		rc <- result
		logutil.Logger.Error(result.Err, zap.String("cluster id", strconv.Itoa(int(ce.ClusterID))))
		return
	}

	ce.HistoryID = his.ID

	// global sync and counter
	wg := &sync.WaitGroup{}
	counter := &ExecutorCounter{}

	// execute check in concurrency
	for _, task := range ce.CheckList {
		ctx := ExecutorContext{
			cluster:   ce,
			checkInfo: task,
			wg:        wg,
			counter:   counter,
		}
		executor := createExecutor(ctx)
		if executor == nil {
			logutil.Logger.Error(fmt.Sprintf("create executor error, invalide source: %s", task.Source), zap.String("cluster id", strconv.Itoa(int(ce.ClusterID))))
			continue
		} else {
			wg.Add(1)
			go executor.ExecuteCheck(rc)
		}
	}
	wg.Wait()
	// update history result
	his.Duration = time.Since(begin).Milliseconds()
	his.NormalItems = counter.normalItems
	his.WarningItems = counter.warningItems
	his.TotalItems = counter.normalItems + counter.warningItems
	his.State = "finished"
	model.DbConn.Save(&his)
	model.DbConn.Model(&model.Cluster{}).Where("id = ?", ce.ClusterID).Update("last_check_time", begin)
	logutil.Logger.Info("cluster check finished.", zap.String("cluster id", strconv.Itoa(int(ce.ClusterID))), zap.Any("result", his))
	// send finish signal
	result := CheckResult{IsFinished: true, CheckID: his.ID}
	rc <- result
	// update clsuter healthy
	go his.UpdateClusterHealthy(int(ce.ClusterID))
}

/*
@clusterID - the cluster id
@schedulerID - the scheduler id, no scheduler trigger use 0
*/
// CreateClusterExecutor create TiDB cluster check executor
func CreateClusterExecutor(clusterID, schedulerID uint) Executor {
	c, err := (&model.Cluster{}).QueryClusterInfoByID(int(clusterID))
	if err != nil {
		logutil.Logger.Error(fmt.Sprintf("CreateClusterExecutor Error: %s", err.Error()), zap.String("cluster id", strconv.Itoa(int(clusterID))))
		return nil
	}
	ce := &ClusterExecutor{
		ClusterID:   clusterID,
		SchedulerID: schedulerID,
		Prometheus:  c.PrometheusURL,
		LoginPath:   c.LoginPath,
	}

	tasks, err := (&model.ClusterChecklist{}).GetEnabledCheckListByClusterID(int(clusterID))
	if err != nil {
		logutil.Logger.Error(fmt.Sprintf("GetClusterChecklist Error: %s", err.Error()), zap.String("cluster id", strconv.Itoa(int(clusterID))))
		return nil
	}

	ce.CheckList = tasks
	return ce
}

// applyProbe use a probe in check
func applyProbe(ctx ExecutorContext, rc chan CheckResult) {
	result := CheckResult{}

	file := fmt.Sprintf("%s/%s/%s/%s", config.GlobalConfig.GetProbePrefix(), ctx.checkInfo.Source, ctx.checkInfo.ProbeID, ctx.checkInfo.FileName)
	_, e := os.Stat(file)
	if os.IsNotExist(e) {
		logutil.Logger.Error(fmt.Sprintf("applyProbe Error, file not found: %s", e.Error()), zap.String("cluster id", strconv.Itoa(int(ctx.cluster.ClusterID))))
		result.Err = e.Error()
		rc <- result
		return
	}

	f, e := filepath.Abs(file)
	if e != nil {
		logutil.Logger.Error(fmt.Sprintf("applyProbe Error, file not found: %s", e.Error()), zap.String("cluster id", strconv.Itoa(int(ctx.cluster.ClusterID))))
		result.Err = e.Error()
		rc <- result
		return
	}
	// the required arguments
	args := []string{f}                         // script file absolute path
	args = append(args, "basepath")             // TODO
	args = append(args, ctx.cluster.LoginPath)  //login path
	args = append(args, ctx.cluster.Prometheus) //promethous url
	args = append(args, ctx.checkInfo.Arg)      //probe custom args

	var output []string
	var err error
	begin := time.Now()
	switch path.Ext(file) {
	case ".sh":
		output, err = applyShellProbe(args)
	case ".py":
		output, err = applyPythonProbe(args)
	default:
		logutil.Logger.Error(fmt.Sprintf("applyProbe Error, invalid file extension: %s", file), zap.String("cluster id", strconv.Itoa(int(ctx.cluster.ClusterID))))
		result.Err = fmt.Sprintf("invalid file extension: %s", path.Ext(file))
		return
	}
	if err != nil {
		cd := model.CheckData{
			Duration:    time.Since(begin).Milliseconds(),
			ProbeID:     ctx.checkInfo.ProbeID,
			CheckTag:    ctx.checkInfo.Tag,
			CheckTime:   begin,
			CheckName:   ctx.checkInfo.ScriptName,
			ClusterID:   ctx.cluster.ClusterID,
			HistoryID:   ctx.cluster.HistoryID,
			CheckStatus: -1,
			CheckItem:   ctx.checkInfo.ScriptName,
			CheckValue:  err.Error(),
		}
		result.Data = append(result.Data, cd)
		model.DbConn.Create(&cd)
	} else {
		normal, warning, data := compareThreshold(ctx, begin, output)
		result.Data = data
		if len(data) > 0 {
			model.DbConn.Create(&data)
		}
		ctx.counter.mutex.Lock()
		ctx.counter.normalItems += normal
		ctx.counter.warningItems += warning
		ctx.counter.mutex.Unlock()
	}
	rc <- result
}

// compareThreshold check probe result
func compareThreshold(
	ctx ExecutorContext,
	begin time.Time,
	output []string) (normal, warning uint, data []model.CheckData) {
	// fmt.Println(len(output))
	duration := time.Since(begin).Milliseconds()
	if len(output) == 0 {
		cd := model.CheckData{
			Duration:  duration,
			ProbeID:   ctx.checkInfo.ProbeID,
			CheckTag:  ctx.checkInfo.Tag,
			CheckTime: begin,
			CheckName: ctx.checkInfo.ScriptName,
			ClusterID: ctx.cluster.ClusterID,
			HistoryID: ctx.cluster.HistoryID,
		}
		cd.Operator = ctx.checkInfo.Operator
		cd.Threshold = ctx.checkInfo.Threshold
		cd.Arg = ctx.checkInfo.Arg
		cd.CheckItem = ctx.checkInfo.ScriptName
		cd.CheckValue = "NA"
		cd.CheckStatus = 0
		data = append(data, cd)
		normal++
		return
	}

	var thd float64
	if t, e := strconv.ParseFloat(ctx.checkInfo.Threshold, 32); e == nil {
		thd = t
	}
	for i := 0; i < len(output); i++ {
		op := output[i]

		cd := model.CheckData{
			Duration:  duration,
			ProbeID:   ctx.checkInfo.ProbeID,
			CheckTag:  ctx.checkInfo.Tag,
			CheckTime: begin,
			CheckName: ctx.checkInfo.ScriptName,
			ClusterID: ctx.cluster.ClusterID,
			HistoryID: ctx.cluster.HistoryID,
		}
		cd.Operator = ctx.checkInfo.Operator
		cd.Threshold = ctx.checkInfo.Threshold
		cd.Arg = ctx.checkInfo.Arg

		is_normal := true
		if strings.HasPrefix(op, "[tck_result:]") {
			row := strings.Split(strings.TrimPrefix(op, "[tck_result:]"), "=")
			if len(row) < 2 {
				logutil.Logger.Warn(fmt.Sprintf("error to skipped: invald prob: %s", op), zap.String("cluster id", strconv.Itoa(int(ctx.cluster.ClusterID))))
				continue
			}
			cd.CheckItem = row[0]
			cd.CheckValue = row[1]

			var val float64
			if t, e := strconv.ParseFloat(row[1], 32); e == nil {
				val = t
			}

			switch ctx.checkInfo.Operator {
			case Comparator_Eq:
				{
					is_normal = val == thd
				}
			case Comparator_Gt:
				{
					is_normal = val > thd
				}
			case Comparator_Ge:
				{
					is_normal = val >= thd
				}
			case Comparator_Le:
				{
					is_normal = val < thd
				}
			case Comparator_Lt:
				{
					is_normal = val <= thd
				}
			case Comparator_NA:
				{
					is_normal = true
				}
			default:
				{
					logutil.Logger.Warn(fmt.Sprintf("error to skipped: comparator %v not supported", ctx.checkInfo.Operator), zap.String("cluster id", strconv.Itoa(int(ctx.cluster.ClusterID))))
					continue
				}
			}

		} else if strings.HasPrefix(op, "[tck_log:]") {
			// TODO: save check log of a script in furtuer
			logutil.Logger.Info(op, zap.String("cluster id", strconv.Itoa(int(ctx.cluster.ClusterID))))

			continue
		} else {
			// the others output as script error to save
			cd.CheckItem = ctx.checkInfo.ScriptName
			cd.CheckValue = op
			is_normal = false
		}
		if is_normal {
			normal = normal + 1
			cd.CheckStatus = 0
		} else {
			warning = warning + 1
			cd.CheckStatus = 1
		}
		data = append(data, cd)
	}
	return
}

// applyShellProbe run a .sh file
func applyShellProbe(args []string) (output []string, err error) {
	cmd := exec.Command("sh", args...)
	op, e := cmd.CombinedOutput()
	if e != nil {
		fmt.Println(e.Error())
		err = e
	} else {
		if s := string(op); len(s) > 0 {
			output = strings.Split(strings.Trim(s, "\n"), "\n")
			//fmt.Println(output)
		}
	}
	return
}

// applyPythonProbe run s python file
func applyPythonProbe(args []string) (output []string, err error) {
	cmd := exec.Command("python3", args...)
	op, e := cmd.CombinedOutput()
	if e != nil {
		fmt.Println(e.Error())
		err = e
	} else {
		if s := string(op); len(s) > 0 {
			output = strings.Split(strings.Trim(s, "\n"), "\n")
			//fmt.Println(output)
		}
	}
	return
}

type ProbeExecutor interface {
	ExecuteCheck(rc chan CheckResult)
}

// type LocalExecutor struct {
// 	LocalExecutorType ExecutorType
// 	Info              *model.CheckListInfo
// 	Cluster           *ClusterExecutor
// }

// func (le *LocalExecutor) ExecuteCheck(res *CheckResult) {
// 	applyProbe(le.Info, le.Cluster, res)
// }

// type RemoteExecutor struct {
// 	RemoteExecutorType ExecutorType
// 	Info               *model.CheckListInfo
// 	Cluster            *ClusterExecutor
// }

// func (re *RemoteExecutor) ExecuteCheck(res *CheckResult) {
// 	applyProbe(re.Info, re.Cluster, res)
// }

// type CustomExecutor struct {
// 	CustomExecutorType ExecutorType
// 	Info               *model.CheckListInfo
// 	Cluster            *ClusterExecutor
// }

// func (ce *CustomExecutor) ExecuteCheck(res *CheckResult) {
// 	applyProbe(ce.Info, ce.Cluster, res)
// }

type commonExecutor struct {
	context ExecutorContext
}

// ExecuteCheck
func (ce *commonExecutor) ExecuteCheck(rc chan CheckResult) {
	applyProbe(ce.context, rc)
	ce.context.wg.Done()
}

func createExecutor(ctx ExecutorContext) ProbeExecutor {
	return &commonExecutor{context: ctx}
	// fmt.Println("Create ProbeExecutor:", info.ScriptName, info.Source)
	// switch info.Source {
	// case string(LocalExecutorType):
	// 	return &LocalExecutor{Info: info, Cluster: cluster}
	// case string(RemoteExecutorType):
	// 	return &RemoteExecutor{Info: info, Cluster: cluster}
	// case string(CustomExecutorType):
	// 	return &CustomExecutor{Info: info, Cluster: cluster}
	// default:
	// 	return nil
	// }
}

type ExecutorType string

const (
	LocalExecutorType  ExecutorType = "local"
	RemoteExecutorType ExecutorType = "remote"
	CustomExecutorType ExecutorType = "custom"
)

const (
	Comparator_NA = iota
	Comparator_Eq
	Comparator_Gt
	Comparator_Ge
	Comparator_Lt
	Comparator_Le
)
