package demo

import (
	"encoding/json"
	"github.com/practice/Golang-Concurrency-Pattern-Demo/task-job-mode/common/constants"
	"github.com/practice/Golang-Concurrency-Pattern-Demo/task-job-mode/model"
	"log"
	"strconv"
	"time"
)

type Plugin struct{}

func (plugin *Plugin) Execute(job *model.Job, sendResponse func(interface{})) (model.ExecStatus, interface{}) {
	jobJson, _ := json.Marshal(job)
	log.Printf("DemoPlugin: received job [%s]", jobJson)
	sendResponse(map[string]interface{}{"message": []string{"Hello!", "central", "operation", "client"}})
	// 模拟执行
	if arg, ok := job.Args["sleepTime"]; ok {
		sleepTime, _ := strconv.Atoi(arg)
		time.Sleep(time.Duration(sleepTime) * 100 * time.Millisecond)
	}
	// 模拟异常
	if _, ok := job.Args["panic"]; ok {
		panic("DemoPlugin: panic")
	}
	data := map[string]interface{}{"message": []string{"DemoPlugin: work complete"}}
	// 模拟执行失败
	if _, ok := job.Args["failure"]; ok {
		return constants.Failure, data
	}
	return constants.Success, data
}

func (plugin *Plugin) ActionList() []*model.Action {
	return []*model.Action{{Name: "demo", Version: "v1.0.0"}}
}
