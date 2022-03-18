package executor

import (
	"TiCheck/internal/model"
	"context"
	"fmt"
)

type CheckResult struct {
	Err  error // script level error
	Data model.CheckData
}

type Executor interface {
	ExecuteCheck(ctx context.Context, result chan CheckResult) // executes one round of check
}

type LocalExecutor struct {
	ID string
}

func (le *LocalExecutor) ExecuteCheck(ctx context.Context, rc chan CheckResult) {
	result := CheckResult{}
	if le.ID == "" {
		result.Err = fmt.Errorf("LocalExecutor: ID is empty")
	} else {

		result.Data = model.CheckData{}
	}
	rc <- result
}

type RemoteExecutor struct {
	ID string
}

func (re *RemoteExecutor) ExecuteCheck(ctx context.Context, rc chan CheckResult) {
	result := CheckResult{}
	result.Err = nil
	result.Data = model.CheckData{}
	rc <- result
}

type CustomExecutor struct {
	ID string
}

func (ce *CustomExecutor) ExecuteCheck(ctx context.Context, rc chan CheckResult) {
	result := CheckResult{}
	result.Err = nil
	result.Data = model.CheckData{}
	rc <- result
}

type ExecutorFactory interface {
	Create(executorType string) Executor
}

/*
@t: executor type
@p: probe id
*/
func Create(t ExecutorType, p string) Executor {
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

type ExecutorType string

const (
	LocalExecutorType  ExecutorType = "local"
	RemoteExecutorType ExecutorType = "remote"
	CustomExecutorType ExecutorType = "custom"
)
