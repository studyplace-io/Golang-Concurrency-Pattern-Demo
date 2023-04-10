package pipeline_mode

import (
	"fmt"
	"testing"
)

const (
	UnExecute = "unexecuted"
	Running   = "running"
	Failed    = "failed"
	Finished  = "finished"
)

type Task struct {
	Id     int
	Action string
	Status string
	Result string
}

func NewTask(action string) *Task {
	task := &Task{
		Id:     1,
		Action: action,
		Status: UnExecute,
	}
	return task
}

// Execute 执行任务
func (t *Task) Execute() *Task {

	// 执行具体任务
	err := t.handler()
	if err != nil {
		t.Status = Failed
		return nil
	}

	// 完成任务
	t.Status = Finished
	return t
}

// handler 处理不同任务
func (t *Task) handler() error {

	// task执行中
	t.Status = Running

	// 暂时这样写，可以调用具体方法
	switch t.Action {
	case "task1":
		fmt.Println("execute task: ", t.Action)
		t.Result = fmt.Sprintf("execute task: %v", t.Action)
	case "task2":
		fmt.Println("execute task: ", t.Action)
		t.Result = fmt.Sprintf("execute task: %v", t.Action)
	case "task3":
		fmt.Println("execute task: ", t.Action)
		t.Result = fmt.Sprintf("execute task: %v", t.Action)
	}

	return nil
}

// PrepareTask 负责准备任务，返回一个chan，把准备好的数据放入chan中。
func PrepareTask(tasks ...*Task) <-chan Task {
	out := make(chan Task)

	// 异步启goroutine准备数据，并放入chan
	go func() {
		defer close(out)
		for _, task := range tasks {
			// 这里可以执行task的预处理
			//
			out <- *task
		}
	}()

	return out
}

// ExecuteTask 执行主要的业务逻辑
func ExecuteTask(inputC <-chan Task) <-chan *Task {

	out := make(chan *Task)
	// 异步启goroutine执行业务逻辑，并放入chan
	go func() {
		defer close(out)
		for task := range inputC {
			out <- task.Execute()
		}
	}()
	return out
}

// AnalyzeTask 消费task任务的结果
func AnalyzeTask(resultTaskC <-chan *Task) {
	for res := range resultTaskC {
		fmt.Println(res.Result)
	}
}

func TestTaskPipeline(t *testing.T) {
	task1 := NewTask("task1")
	task2 := NewTask("task2")
	task3 := NewTask("task3")

	// 流水线
	prepareTaskC := PrepareTask(task3, task2, task1)
	resultTaskC := ExecuteTask(prepareTaskC)
	AnalyzeTask(resultTaskC)

}
