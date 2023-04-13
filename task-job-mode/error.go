package task_job_mode

import "errors"

var AlreadyClosed = errors.New("worker already closed")

type ErrorResult struct {
	ErrorNo  int    `json:"errorNo"`
	ErrorMsg string `json:"errorMsg"`
}

var (
	PluginExecuteError = &ErrorResult{ErrorNo: 10000, ErrorMsg: "plugin execute job error"}
	NoPluginHandle     = &ErrorResult{ErrorNo: 10001, ErrorMsg: "no plugin can handle this job"}
)
