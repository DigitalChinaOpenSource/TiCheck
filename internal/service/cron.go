package service

import (
	"TiCheck/executor"
	"TiCheck/internal/model"
	"github.com/robfig/cron/v3"
	"time"
)

type SchedulerCron struct {
	Cron  *cron.Cron
	Tasks []cron.EntryID
}

type SchedulerTask struct {
	ID          cron.EntryID
	SchedulerID uint
	ClusterID   uint
	Spec        string
	Status      bool
	JobFunc     func()
}

var CronService SchedulerCron

func (st *SchedulerTask) Initialize() {
	CronService.Cron = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
	CronService.Cron.Start()
}

func (st *SchedulerTask) AddTask(scheduler model.Scheduler) {
	st.Spec = scheduler.CronExpression
	st.CreateJob(scheduler)
	//
	//CronService.Cron.AddFunc()
}

func (st *SchedulerTask) CreateJob(scheduler model.Scheduler) {
	jobFunc := func() {
		exe := executor.CreateClusterExecutor(scheduler.ClusterID, scheduler.ID)
		resultCh := make(chan executor.CheckResult, 10)
		// ctx := context.WithValue(context.Background(), "", "")
		go exe.Execute(resultCh)
	}
	st.JobFunc = jobFunc
}
