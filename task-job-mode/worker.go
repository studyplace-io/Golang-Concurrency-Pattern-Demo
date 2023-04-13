package task_job_mode

import (
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/common/constants"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/model"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/plugins"
	"log"
)

// executeTask 执行任务
func executeTask(newTask *taskNode) {
	taskNode := newTask
	// 退出时减少并发数
	defer func() {
		<-pool.workerChan
	}()
	var ok bool
	// 不断拉取任务执行
	for taskNode != nil {
		for index := range taskNode.task.Jobs {
			executeJob(taskNode, index)
		}
		//completeTask(taskNode, constants.Completed)
		// 拉取新任务
		if taskNode, ok = pool.queue.requestTask(); !ok {
			break
		}
	}
}

// executeJob 执行小任务
func executeJob(taskNode *taskNode, jobIndex int) {
	// job 信息
	job := &taskNode.task.Jobs[jobIndex]
	status, taskStatus := constants.Running, constants.Running
	// 提供给plugin的，用于向socket响应数据用
	responseFunc := func(data interface{}) {
		result := &model.Result{
			JobId:      job.JobId,
			TaskId:     taskNode.task.TaskId,
			Status:     string(status),
			TaskStatus: string(taskStatus),
			Data:       data,
		}
		//_ = dal.CreateJobResult(nil, taskNode.jobIds[jobIndex], data)
		if err := pool.responseWriter(result); err != nil {
			log.Printf("Send result to server failed: %v", err)
		}
	}
	// 执行结果回复
	var data interface{}
	defer func() {
		// 执行异常时恢复
		if err := recover(); err != nil {
			log.Printf("Execute task [%#v] Job [%#v] panic: %v", taskNode.task, job, err)
			status, data = constants.Failure, PluginExecuteError
		}
		// 最后一个Job完成时，标记Task完成
		if jobIndex == len(taskNode.task.Jobs)-1 {
			taskStatus = constants.Completed
		}
		// 回复Server结果
		responseFunc(data)
		// 标记任务状态
		//completeJob(taskNode.jobIds[jobIndex], status)
	}()
	// 调用插件执行
	if plugin, ok := plugins.SearchPlugin(job.Service, job.Action); ok {
		status, data = plugin.Execute(job, responseFunc)
		log.Printf("Job [%s] done, completion state [%s]", job.JobId, status)
	} else {
		log.Printf("No plugin support for job [%s]: service [%s] action [%s]", job.JobId, job.Service, job.Action)
		status, data = constants.Failure, NoPluginHandle
	}
}

//func completeTask(task *taskNode, status model.ExecStatus) {
// // 将任务执行状态同步到数据库,并标记为已完成
// doneAt := time.Now()
// _ = dal.UpdateTaskStatusById(nil, task.id, status, &doneAt)
//}
//
//func completeJob(jobId int64, status model.ExecStatus) {
// // 将Job完成状态记录至数据库
// doneAt := time.Now()
// _ = dal.UpdateJobStatusById(nil, jobId, status, &doneAt)
//}
