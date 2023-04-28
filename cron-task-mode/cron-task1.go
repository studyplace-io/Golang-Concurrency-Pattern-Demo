package cron_practice

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

/*
	参考：
	https://juejin.cn/post/7004656484902502408
	http://www.zyiz.net/tech/detail-141215.html
*/

func TestCron(wg *sync.WaitGroup) {

	defer wg.Done()

	c := cron.New()

	i := 1
	EntryID, err := c.AddFunc("*/1 * * * *", func() {
		fmt.Println(time.Now(), "每分钟执行一次", i)
	})
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
	//time.Sleep(time.Minute * 1)

}

func CronTask1() {

	var wg sync.WaitGroup

	wg.Add(2)
	go TestCron(&wg)
	go TestCron2(&wg)

	wg.Wait()
}

type Job1 struct {
}

func (j *Job1) Run() {
	fmt.Println(time.Now(), "Job1开始工作")
}

type Job2 struct {
}

func (j *Job2) Run() {
	fmt.Println(time.Now(), "Job2开始工作")
}

func TestCron2(wg *sync.WaitGroup) {

	defer wg.Done()

	c := cron.New(cron.WithSeconds())

	EntryID, err := c.AddJob("*/3 * * * * *", &Job1{})
	fmt.Println(time.Now(), EntryID, err)

	EntryID, err = c.AddJob("*/5 * * * * *", &Job2{})
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
	time.Sleep(time.Second * 10)

}
