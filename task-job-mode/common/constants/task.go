package constants

import "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/model"

// 状态。目前task、job及plugin共用。
const (
	Running   model.ExecStatus = "running"
	Success   model.ExecStatus = "success"   // 成功
	Failure   model.ExecStatus = "failure"   // 失败
	Aborted   model.ExecStatus = "aborted"   // 取消
	Completed model.ExecStatus = "completed" // 完成，仅task
)
