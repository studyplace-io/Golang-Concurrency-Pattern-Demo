package worker

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"sync"
	"time"
)

// Worker 接口
type Worker interface {
	// RegisterJob 注册 Job 任务，输入 Job name
	RegisterJob(name string)
	// RunCronJob 定时执行 Job ，输入轮循时间、任务逻辑 func
	RunCronJob(name string, period time.Duration, f func()) error
	// RunCronJobWithContext 定时执行 Job 并传入 context
	RunCronJobWithContext(name string, ctx context.Context,
		period time.Duration, f func(ctx context.Context)) error
	// IsExist 是否存在
	IsExist(name string) bool
	// JobStatus Job 状态
	JobStatus(name string) string
	// StopJob 停止 Job
	StopJob(name string)
	// StopAll 停止全部 Job
	StopAll()
	// Range 遍历
	Range()
}

const (
	Running     = "running"
	NotExecuted = "notexecuted"
	Stopped     = "stopped"
)

// worker 实现 Worker 对象
type worker struct {
	// JobMap 存储 Job 任务
	JobMap sync.Map
}

// job 任务
type job struct {
	// job 名
	name string
	// status 状态
	status string
	// 停止时使用
	stopCh chan struct{}
}

func NewWorker() Worker {
	return &worker{
		JobMap: sync.Map{},
	}
}

func (w *worker) IsExist(name string) bool {
	_, ok := w.JobMap.Load(name)
	return ok
}

func (w *worker) JobStatus(name string) string {
	j, ok := w.JobMap.Load(name)
	if !ok {
		return ""
	}
	jobObj := j.(*job)
	return jobObj.status
}

func (w *worker) RegisterJob(name string) {
	_, ok := w.JobMap.Load(name)
	if !ok {
		task := &job{
			name:   name,
			status: NotExecuted,
			stopCh: make(chan struct{}),
		}
		w.JobMap.Store(name, task)
	}
}

func (w *worker) RunCronJob(name string, period time.Duration, f func()) error {
	j, ok := w.JobMap.Load(name)
	if !ok {
		return fmt.Errorf("%s not registered", name)
	}

	jobObj := j.(*job)
	jobObj.status = Running
	go wait.Until(f, period, jobObj.stopCh)
	return nil
}

func (w *worker) RunCronJobWithContext(name string, ctx context.Context, period time.Duration, f func(ctx context.Context)) error {
	j, ok := w.JobMap.Load(name)
	if !ok {
		return fmt.Errorf("%s not registered", name)
	}

	jobObj := j.(*job)
	jobObj.status = Running
	go wait.JitterUntilWithContext(ctx, f, period, 0, true)
	return nil
}

func (w *worker) StopJob(name string) {
	j, ok := w.JobMap.Load(name)
	if !ok {
		return
	}

	jobObj := j.(*job)
	jobObj.status = Stopped
	close(jobObj.stopCh)
	w.JobMap.Delete(name)
}

func (w *worker) Range() {
	w.JobMap.Range(func(key, value any) bool {
		fmt.Printf("worker name: [%s]\n", key)
		return true
	})
}

func (w *worker) StopAll() {
	w.JobMap.Range(func(key, value interface{}) bool {
		jobObj := value.(*job)
		jobObj.status = Stopped
		close(jobObj.stopCh)
		w.JobMap.Delete(key)
		return true
	})
}
