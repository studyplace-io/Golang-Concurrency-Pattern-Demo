package task_job_mode

import (
	"encoding/json"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/common/constants"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/model"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestQueue 测试任务队列基本功能
func TestQueue(t *testing.T) {
	Convey("Queue Test", t, func() {
		taskNodes := make([]*taskNode, 100)
		var wg sync.WaitGroup
		q := &taskQueue{}
		// push
		go func() {
			wg.Add(1)
			defer wg.Done()
			for i := 0; i < 100; i++ {
				taskNodes[i] = &taskNode{id: int64(i), task: &model.Task{TaskId: strconv.Itoa(i)}}
				q.push(taskNodes[i])
			}
		}()
		// poll
		for i := 0; i < 95; {
			if task, ok := q.requestTask(); ok {
				So(task, ShouldNotBeNil)
				i++
				taskNodes[task.id].id = -1
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
		wg.Wait()
		abortCount := 0
		So(q.length(), ShouldEqual, 5)
		// 测试任务拦截
		for _, node := range taskNodes {
			if node.id != -1 {
				ok := q.abortTask(node.task.TaskId)
				So(ok, ShouldBeTrue)
				abortCount++
			}
		}
		So(abortCount, ShouldEqual, 5)
		// 测试已分配任务拦截
		for i := 0; i < 10; i++ {
			So(q.abortTask(taskNodes[i].task.TaskId), ShouldBeFalse)
		}
		So(q.length(), ShouldEqual, 0)
	})
}

func newTestUuid() string {
	return uuid.New().String()
}

func writer() func(result *model.Result) error {
	return func(result *model.Result) error {
		resp, _ := json.Marshal(result)
		log.Printf("Response: %s", resp)
		return nil
	}
}

// TestClientWorker 测试ClientWorker执行流程
func TestClientWorker(t *testing.T) {
	Convey("CreateClientWorker", t, func() {
		worker := CreateClientWorker(writer())
		So(worker, ShouldNotBeNil)
		So(worker, ShouldEqual, CreateClientWorker(writer()))
	})
	defer func() {
		_ = pool.Close()
	}()
	Convey("ExecTask", t, func() {
		for i := 0; i < 100; i++ {
			err := pool.Execute(&model.Task{
				TaskId: newTestUuid(),
				Jobs: []model.Job{
					{
						JobId:   newTestUuid(),
						Service: "demo",
						Action:  "demo",
						Args:    map[string]string{"sleepTime": strconv.Itoa(i%5 + 1)},
					},
				},
			})

			So(err, ShouldBeNil)
		}
	})
}

func TestAbnormalExecute(t *testing.T) {
	// 异常场景测试
	Convey("TestExecuteAbnormalScene", t, func() {
		var panics, noPlugin, failure int
		_ = CreateClientWorker(func(result *model.Result) error {
			resp, _ := json.Marshal(result)
			log.Printf("Response: %s", resp)
			if err, ok := result.Data.(*ErrorResult); ok {
				switch err.ErrorNo {
				case 10000:
					panics++
				case 10001:
					noPlugin++
				}
			}
			if result.Status == string(constants.Failure) {
				failure++
			}
			return nil
		})
		// 测试无plugin可处理请求
		noPluginTimes := 5
		for i := 0; i < noPluginTimes; i++ {
			err := pool.Execute(&model.Task{
				TaskId: newTestUuid(),
				Jobs: []model.Job{
					{
						JobId:   newTestUuid(),
						Service: "demo",
						Action:  "demo1",
					},
				},
			})
			So(err, ShouldBeNil)
		}
		// 测试plugin panic
		panicTimes := 4
		for i := 0; i < panicTimes; i++ {
			err := pool.Execute(&model.Task{
				TaskId: newTestUuid(),
				Jobs: []model.Job{
					{
						JobId:   newTestUuid(),
						Service: "demo",
						Action:  "demo",
						Args:    map[string]string{"panic": "1"},
					},
				},
			})
			So(err, ShouldBeNil)
		}

		// 测试任务执行失败
		failedTimes := 6
		for i := 0; i < failedTimes; i++ {
			err := pool.Execute(&model.Task{
				TaskId: newTestUuid(),
				Jobs: []model.Job{
					{
						JobId:   newTestUuid(),
						Service: "demo",
						Action:  "demo",
						Args:    map[string]string{"failure": "1"},
					},
				},
			})
			So(err, ShouldBeNil)
		}

		_ = pool.Close()

		So(noPlugin, ShouldEqual, noPluginTimes)
		So(panics, ShouldEqual, panicTimes)
		So(failure, ShouldEqual, failedTimes+panicTimes+noPluginTimes)
	})

	// 重复执行测试
	Convey("TestExecuteRepeatTask", t, func() {
		worker := CreateClientWorker(writer())
		defer worker.Close()
		task := &model.Task{
			TaskId: newTestUuid(),
			Jobs: []model.Job{
				{
					JobId:   newTestUuid(),
					Service: "demo",
					Action:  "demo",
				},
			},
		}
		err := worker.Execute(task)
		So(err, ShouldBeNil)
		err = worker.Execute(task)
		So(err, ShouldBeNil)
	})
}

func TestExecuteAfterClose(t *testing.T) {
	// 异常场景测试
	Convey("TestExecuteAfterClose", t, func() {
		// 测试异常配置
		//common.AppConfig.Worker.ConcurrencyLimit = -1
		worker := CreateClientWorker(writer())
		err := worker.Execute(&model.Task{
			TaskId: newTestUuid(),
			Jobs: []model.Job{
				{
					JobId:   newTestUuid(),
					Service: "demo",
					Action:  "demo",
				},
			},
		})
		So(err, ShouldBeNil)
		err = worker.Close()
		So(err, ShouldBeNil)
		// 测试重复关闭
		err = worker.Close()
		So(err, ShouldBeNil)
		// 关闭后发送
		err = worker.Execute(&model.Task{
			TaskId: newTestUuid(),
			Jobs: []model.Job{
				{
					JobId:   newTestUuid(),
					Service: "demo",
					Action:  "demo",
				},
			},
		})
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, AlreadyClosed)
	})
}

// TestBatchJobs 测试Task下多Job的情况
func TestBatchJobs(t *testing.T) {
	Convey("TestBatchJobs", t, func() {
		tasks, jobs := 10, 30
		var completedTask, completedJob int
		worker := CreateClientWorker(func(result *model.Result) error {
			resp, _ := json.Marshal(result)
			if result.TaskStatus == string(constants.Completed) {
				completedTask++
			}
			if result.Status == string(constants.Success) || result.Status == string(constants.Failure) {
				completedJob++
			}
			log.Printf("Response: %s", resp)
			return nil
		})
		for i := 0; i < tasks; i++ {
			err := worker.Execute(&model.Task{
				TaskId: newTestUuid(),
				Jobs: []model.Job{
					{
						JobId:   newTestUuid(),
						Service: "demo",
						Action:  "demo",
					},
					{
						JobId:   newTestUuid(),
						Service: "demo",
						Action:  "demo",
					},
					{
						JobId:   newTestUuid(),
						Service: "demo",
						Action:  "demo",
					},
				},
			})

			So(err, ShouldBeNil)
		}
		_ = worker.Close()
		So(tasks, ShouldEqual, completedTask)
		So(jobs, ShouldEqual, completedJob)
	})
}
