package executor

import (
	"TiCheck/internal/model"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

var (
	probe_prefix = "../probes"
)

type CheckResult struct {
	Err  error // script level error
	Data model.CheckData
}

type Executor interface {
	Execute(result chan CheckResult) // executes one round of check
}

type ClusterExecutor struct {
	ClusterID  int
	Prometheus string
	LoginPath  string
	CheckList  []model.CheckListInfo
}

func (ce *ClusterExecutor) Execute(rc chan CheckResult) {

	for _, task := range ce.CheckList {

		result := CheckResult{}
		executor := createExecutor(&task, ce)
		if executor == nil {
			result.Err = fmt.Errorf("create executor error, invalide source: %s", task.Source)
			continue
		} else {
			go executor.ExecuteCheck(result)
		}
		rc <- result
	}
}

func CreateClusterExecutor(cluster_id int) Executor {
	c, err := (&model.Cluster{}).QueryClusterInfoByID(cluster_id)
	if err != nil {
		fmt.Println("CreateClusterExecutor Error:", err.Error())
		return nil
	}
	ce := &ClusterExecutor{
		ClusterID:  cluster_id,
		Prometheus: c.PrometheusURL,
		LoginPath:  c.LoginPath,
	}

	tasks, err := (&model.ClusterChecklist{}).GetListInfoByClusterID(cluster_id)
	if err != nil {
		fmt.Println("CreateClusterExecutor Error:", err.Error())
		return nil
	}

	ce.CheckList = tasks
	return ce
}

func applyProbe(file string, info *model.CheckListInfo, cluster *ClusterExecutor) CheckResult {
	result := CheckResult{}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Println("applyProbe Error, file not found:", err.Error())
		result.Err = err
		return result
	}

	args := []string{file}                  // sh file
	args = append(args, cluster.Prometheus) //promethous url
	args = append(args, cluster.LoginPath)  //login path
	args = append(args, info.Arg)           //probe custom args

	switch path.Ext(file) {
	case ".sh":
		result = applyShellProbe(args)
	case ".py":
		result = applyPythonProbe(args)
	default:
		fmt.Println("applyProbe Error, invalid file extension:", file)
		result.Err = fmt.Errorf("invalid file extension: %s", path.Ext(file))
	}

	return result
}

func applyShellProbe(args []string) CheckResult {
	cmd := exec.Command("sh", args...)
	err := cmd.Run()
	if err != nil {
		print(err)
	}
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	// sleep for 10 seconds before sending done signal
	time.Sleep(time.Second * 10)
	return CheckResult{}
}

func applyPythonProbe(args []string) CheckResult {
	cmd := exec.Command("python", args...)
	err := cmd.Run()
	if err != nil {
		print(err)
	}

	// sleep for 10 seconds before sending done signal
	time.Sleep(time.Second * 10)
	return CheckResult{}
}

type ProbeExecutor interface {
	ExecuteCheck(res CheckResult)
}

type LocalExecutor struct {
	LocalExecutorType ExecutorType
	Info              *model.CheckListInfo
	Cluster           *ClusterExecutor
}

func (le *LocalExecutor) ExecuteCheck(res CheckResult) {

	file := fmt.Sprintf("%s/%s/%s/%s", probe_prefix, le.LocalExecutorType, le.Info.ProbeID, le.Info.FileName)
	res = applyProbe(file, le.Info, le.Cluster)
}

type RemoteExecutor struct {
	RemoteExecutorType ExecutorType
	Info               *model.CheckListInfo
	Cluster            *ClusterExecutor
}

func (re *RemoteExecutor) ExecuteCheck(res CheckResult) {
	file := fmt.Sprintf("%s/%s/%s/%s", probe_prefix, re.RemoteExecutorType, re.Info.ProbeID, re.Info.FileName)
	res = applyProbe(file, re.Info, re.Cluster)
}

type CustomExecutor struct {
	CustomExecutorType ExecutorType
	Info               *model.CheckListInfo
	Cluster            *ClusterExecutor
}

func (ce *CustomExecutor) ExecuteCheck(res CheckResult) {
	file := fmt.Sprintf("%s/%s/%s/%s", probe_prefix, ce.CustomExecutorType, ce.Info.ProbeID, ce.Info.FileName)
	res = applyProbe(file, ce.Info, ce.Cluster)
}

func createExecutor(info *model.CheckListInfo, cluster *ClusterExecutor) ProbeExecutor {
	fmt.Println("Create ProbeExecutor:", info.ScriptName, info.Source)
	switch info.Source {
	case string(LocalExecutorType):
		return &LocalExecutor{Info: info, Cluster: cluster}
	case string(RemoteExecutorType):
		return &RemoteExecutor{Info: info, Cluster: cluster}
	case string(CustomExecutorType):
		return &CustomExecutor{Info: info, Cluster: cluster}
	default:
		return nil
	}
}

type ExecutorType string

const (
	LocalExecutorType  ExecutorType = "local"
	RemoteExecutorType ExecutorType = "remote"
	CustomExecutorType ExecutorType = "custom"
)
