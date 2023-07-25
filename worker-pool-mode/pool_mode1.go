package worker_pool_mode

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Task接口，需要实现Execute方法
type Task interface {
	Execute()
}

// Pool 工作池
type Pool struct {
	TaskC chan Task // 存储任务
	Wg    sync.WaitGroup
	size  int // 任务大小
}

// NewPool 创建pool池，可以传入并发的数量，默认使用cpu数量
func NewPool(size ...int) *Pool {

	if size[0] == 0 {
		size[0] = runtime.NumCPU()
	}

	p := &Pool{
		TaskC: make(chan Task, size[0]),
		Wg:    sync.WaitGroup{},
		size:  size[0],
	}

	return p
}

// Start 开始
func (p *Pool) Start() {
	for i := 0; i < p.size; i++ {
		go func() {
			for task := range p.TaskC {
				task.Execute()
			}
		}()
	}
}

// Stop 停止
func (p *Pool) Stop() {
	defer p.Wg.Wait()
}

func (p *Pool) AddTask(t Task) {
	p.Wg.Add(1)
	p.TaskC <- t
}

type FuncJob func()

func (f FuncJob) Execute() { f() }

func (p *Pool) AddTaskFunc(f FuncJob) {
	p.Wg.Add(1)
	p.TaskC <- f
}

// Task1 任务1
type Task1 struct {
	wg *sync.WaitGroup
}

func (t *Task1) Execute() {
	defer t.wg.Done()
	fmt.Println("Task1 cost 3 seconds")
	time.Sleep(3 * time.Second)
	fmt.Println("Task1 finished!")
}

// Task2 任务2
type Task2 struct {
	wg *sync.WaitGroup
}

func (t *Task2) Execute() {
	defer t.wg.Done()
	fmt.Println("Task2 cost 3 seconds")
	time.Sleep(3 * time.Second)
	fmt.Println("Task2 finished!")
}

// Task3 任务3
type Task3 struct {
	wg *sync.WaitGroup
}

func (t *Task3) Execute() {
	defer t.wg.Done()
	fmt.Println("Task3 cost 3 seconds")
	time.Sleep(3 * time.Second)
	fmt.Println("Task3 finished!")
}
