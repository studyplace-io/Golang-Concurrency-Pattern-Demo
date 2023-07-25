package worker_pool_mode

import (
	"fmt"
	"testing"
	"time"
)


func TestTaskPool(t *testing.T) {
	p := NewPool(2)
	p.Start() // 启动任务

	task1 := &Task1{
		wg: &p.Wg,
	}
	task2 := &Task2{
		wg: &p.Wg,
	}

	task3 := &Task3{
		wg: &p.Wg,
	}

	// 加入任务
	p.AddTask(task1)
	p.AddTask(task2)
	p.AddTask(task3)
	p.AddTaskFunc(func() {
		defer p.Wg.Done()
		fmt.Println("这是一个func!!")
		time.Sleep(time.Second * 5)
		fmt.Println("func finished")
	})
	p.AddTask(task3)

	// 等待所有任务都执行完成
	p.Stop()

}
