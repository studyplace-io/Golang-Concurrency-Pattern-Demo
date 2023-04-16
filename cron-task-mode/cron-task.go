package cron_practice

import (
	"github.com/antlabs/timer"
	"log"
	"sync"
	"time"
)

/*
	使用"github.com/antlabs/timer" 实现定时的简单任务调度
 */


func After(tm timer.Timer, duration time.Duration, callbackFunc func(), taskNum int) {
	var wg sync.WaitGroup
	wg.Add(taskNum)
	for i := 0; i < taskNum; i++ {

		go func() {
			defer wg.Done()
			tm.AfterFunc(duration, callbackFunc)
		}()

	}


	wg.Wait()

}


// ＠param timer.Timer 计时器
// ＠param time.Duration 调度时间
// ＠param op func() 待调度的callback函数
func Schedule(timer timer.Timer, duration time.Duration, callbackFunc func()) {
	timer.ScheduleFunc(duration, callbackFunc)
}


func CronTask() {

	tm := timer.NewTimer()
	defer tm.Stop()


	aa := func() {
		log.Printf("after 1 second stopped!!!\n")
	}
	// 启动一次性任务
	go After(tm, time.Second, aa, 5)


	bb := func() {
		log.Printf("定时调用任务！！\n")
	}
	t := time.Second


	// 启动定时任务
	go Schedule(tm, t, bb)

	go func() {
		time.Sleep(2*time.Minute + 50*time.Second)
		tm.Stop()
	}()

	tm.Run()


}
