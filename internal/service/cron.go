package service

import (
	"TiCheck/executor"
	"TiCheck/internal/model"
	"github.com/robfig/cron/v3"
	"time"
)

var CronService SchedulerCron

type SchedulerCron struct {
	Cron  *cron.Cron
	Tasks []SchedulerTask
}

type SchedulerTask struct {
	ID          cron.EntryID
	SchedulerID uint
	Spec        string
	JobFunc     func()
}

// Initialize initial a cron service
func (sc *SchedulerCron) Initialize() {
	sc.Cron = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
	sc.Cron.Start()
	// todo check if there is a task in the scheduler table, if so,add it into the cron task queue
}

//AddTask add a cron task to cron service by scheduler
func (sc *SchedulerCron) AddTask(scheduler model.Scheduler) error {
	job := CreateJob(scheduler)
	taskID, err := sc.Cron.AddFunc(scheduler.CronExpression, job)
	if err != nil {
		return err
	}
	var task = SchedulerTask{
		taskID,
		scheduler.ID,
		scheduler.CronExpression,
		job,
	}
	sc.Tasks = append(sc.Tasks, task)
	return nil
}

// RemoveTask remove a task from cron service by task info
func (sc *SchedulerCron) RemoveTask(task SchedulerTask) {
	sc.Cron.Remove(task.ID)
}

func (sc *SchedulerCron) StopAll() {
	sc.Cron.Stop()
}

// CreateJob create a cron job by scheduler cron expression, return a job func
func CreateJob(scheduler model.Scheduler) func() {
	jobFunc := func() {
		exe := executor.CreateClusterExecutor(scheduler.ClusterID, scheduler.ID)
		resultCh := make(chan executor.CheckResult, 10)
		// ctx := context.WithValue(context.Background(), "", "")
		go exe.Execute(resultCh)
	}
	return jobFunc
}
