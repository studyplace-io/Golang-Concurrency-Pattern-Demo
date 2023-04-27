package timewheel

import (
	"fmt"
	"testing"
	"time"
)

/*
	参考：https://lk668.github.io/2021/04/05/2021-04-05-%E6%89%8B%E6%8A%8A%E6%89%8B%E6%95%99%E4%BD%A0%E5%A6%82%E4%BD%95%E7%94%A8golang%E5%AE%9E%E7%8E%B0%E4%B8%80%E4%B8%AAtimewheel/
*/

func TestTimeWheel(test *testing.T) {

	// 初始化时间间隔是1s，一共有60个齿轮的时间轮盘，默认轮盘转动一圈的时间是60s
	tw := CreateTimeWheel(1*time.Second, 60, DefaultJob)

	// 启动时间轮
	tw.Start()

	// 关闭时间轮
	defer func() {
		tw.Stop()
	}()

	if tw.IsRunning() {
		// 添加一个task
		// task名字叫task1
		// task的创建时间是time.Now()
		// task执行的任务设置为nil，所以默认执行timewheel的Job，也就是TimeWheelDefaultJob
		fmt.Println(fmt.Sprintf("%v Add task task-5s", time.Now().Format(time.RFC3339)))
		err := tw.AddTask(5*time.Second, "task-5s", time.Now(), -1, nil)
		if err != nil {
			panic(err)
		}

		// 该Task执行example.TaskJob
		fmt.Println(fmt.Sprintf("%v Add task task-2s", time.Now().Format(time.RFC3339)))
		err = tw.AddTask(2*time.Second, "task-2s", time.Now(), -1, TaskJob)
		if err != nil {
			panic(err)
		}

	} else {
		panic("TimeWheel is not running")
	}
	time.Sleep(10 * time.Second)

	// 删除task
	fmt.Println("Remove task task-5s")
	err := tw.RemoveTask("task-5s")
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	fmt.Println("Remove task task-2s")
	err = tw.RemoveTask("task-2s")
	if err != nil {
		panic(err)
	}

}

func DefaultJob(taskName TaskName) {
	fmt.Println(fmt.Sprintf("%v This is a timewheel job with key: %v", time.Now().Format(time.RFC3339), taskName))
}

func TaskJob(taskName TaskName) {
	fmt.Println(fmt.Sprintf("%v This is a task job with key: %v", time.Now().Format(time.RFC3339), taskName))
}
