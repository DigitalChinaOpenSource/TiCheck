package executor

import (
	"TiCheck/internal/model"
	"fmt"
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
	CheckList  []model.Probe
}

func (ce *ClusterExecutor) Execute(result chan CheckResult) {

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

	return ce
}

func runExecutorTask(t ExecutorType, tasks []model.Probe, resultCh chan CheckResult) {

	prefix := fmt.Sprintf("../probes/%s", t)
	for _, task := range tasks {
		file := fmt.Sprintf("%s/%s/%s", prefix, task.ID, task.FileName)
		fmt.Println("runExecutorTask:", file)
		result := applyProbe(file)
		resultCh <- result
	}
}

func applyProbe(file string) CheckResult {
	result := CheckResult{}
	return result
}

type ProbeExecutor interface {
	ExecuteCheck() CheckResult
}

type LocalExecutor struct {
	LocalExecutorType ExecutorType
	ID                string
}

func (le *LocalExecutor) ExecuteCheck() CheckResult {
	result := CheckResult{}
	if le.ID == "" {
		result.Err = fmt.Errorf("LocalExecutor: ID is empty")
	} else {

		result.Data = model.CheckData{}
	}
	return result
}

type RemoteExecutor struct {
	RemoteExecutorType ExecutorType
	ID                 string
}

func (re *RemoteExecutor) ExecuteCheck() CheckResult {
	result := CheckResult{}
	result.Err = nil
	result.Data = model.CheckData{}
	return result
}

type CustomExecutor struct {
	CustomExecutorType ExecutorType
	ID                 string
}

func (ce *CustomExecutor) ExecuteCheck() CheckResult {
	result := CheckResult{}
	result.Err = nil
	result.Data = model.CheckData{}
	return result
}

/*
@t: executor type
@p: probe id
*/
func Create(t ExecutorType, p string) ProbeExecutor {
	fmt.Println("Create executor:", t)
	switch t {
	case LocalExecutorType:
		return &LocalExecutor{ID: p}
	case RemoteExecutorType:
		return &RemoteExecutor{ID: p}
	case CustomExecutorType:
		return &CustomExecutor{ID: p}
	default:
		return &LocalExecutor{ID: p}
	}
}

type ExecutorType string

const (
	LocalExecutorType  ExecutorType = "local"
	RemoteExecutorType ExecutorType = "remote"
	CustomExecutorType ExecutorType = "custom"
)
