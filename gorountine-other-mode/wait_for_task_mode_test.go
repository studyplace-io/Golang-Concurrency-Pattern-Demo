package gorountine_other_mode

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// TestWaitForTResult
// 1.等待结果模式
// 使用无缓存chan实现两个goroutine的通知机制。
func TestWaitForTResult(t *testing.T) {
	waitC := make(chan struct{})

	go func() {
		// 执行业务逻辑
		time.Sleep(time.Second * 5)
		fmt.Println("do something.....")

		// 当执行完后，可以通知主goroutine执行后面的逻辑
		waitC <- struct{}{}
	}()

	fmt.Println("do something other.....")

	<-waitC
	fmt.Println("接著干剩下的活")

}

// Task 接口对象
type Task interface {
	Execute(duration time.Duration)
}

var _ Task = &task{}

type task struct {
	id          int
	status      string
	executeTime time.Duration
}

func newTask(id int) *task {
	return &task{id: id, status: "unrun"}
}

func (t *task) Execute(duration time.Duration) {
	t.executeTime = duration
	t.status = "running"
	time.Sleep(duration)
	fmt.Println("task ", t.id, "is finished")
	t.status = "finished"
}

// TestWaitForTaskOne 等待任务模式
// 使用无缓冲chan+子goroutine通过chan接收主goroutine发送的数据，也可以是执行任务的函数。
func TestWaitForTaskOne(t *testing.T) {
	taskC := make(chan Task)
	var wg sync.WaitGroup

	// 启异步goroutine执行task任务
	wg.Add(1)
	go func(taskC chan Task) {
		defer wg.Done()
		t := <-taskC
		t.Execute(time.Second)

	}(taskC)

	// 业务逻辑
	time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)

	// 主goroutine发送任务
	tt := newTask(10)
	taskC <- tt

	// 执行业务逻辑

	wg.Wait()
}

// TestWaitForTasks 等待任务模式
// 有缓冲chan+子goroutine通过chan接收来自主goroutine发送的数据，也可以是执行任务的函数。
func TestWaitForTasks(t *testing.T) {
	taskC := make(chan Task, 2)
	stopC := make(chan struct{})
	var wg sync.WaitGroup

	// 启异步goroutine执行task任务
	wg.Add(1)
	go func(taskC chan Task) {
		defer wg.Done()
		for {
			select {
			case t := <-taskC:
				t.Execute(time.Second)
			case <-stopC:
				fmt.Println("task execute stopped!")
				return
			}

		}

	}(taskC)

	// 业务逻辑

	// 主goroutine发送任务
	for i := 0; i < 20; i++ {
		taskC <- newTask(rand.Intn(20))
	}

	close(stopC)

	// 执行业务逻辑

	wg.Wait()

}

const (
	capacity   = 10
	requestNum = 30
)

// TestDropMode Drop模式
// 当写入的数据量大的时候，超出chan的容量就选择"随机"丢弃数据或不执行。
// ex: 当负载均衡，或是请求量过大时，可以随机丢弃一些请求。
func TestDropMode(test *testing.T) {
	taskC := make(chan Task, capacity)

	// 异步执行task，也可以使用工作池来并发处理。
	var wg sync.WaitGroup
	wg.Add(1)
	go func(taskC chan Task) {
		defer wg.Done()
		for task := range taskC {
			fmt.Println("execute task")
			task.Execute(time.Second)
		}
	}(taskC)

	// 执行多个任务，且随机丢弃或不执行
	for i := 0; i < requestNum; i++ {
		select {
		case taskC <- newTask(rand.Intn(20)):
			fmt.Println("put task in taskC")
		default:
			fmt.Println("drop some task")
			time.Sleep(time.Second * 1)
		}
	}
	// 关闭
	defer close(taskC)
	wg.Done()

}
