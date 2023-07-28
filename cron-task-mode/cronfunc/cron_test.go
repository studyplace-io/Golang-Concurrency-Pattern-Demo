package cronfunc

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func TestUntil(test *testing.T) {
	// 测试用例
	CronTry()
	CronTry1()
	//CronTry2()
	//CronTry3()
}

func CronTry() {
	ch := make(chan struct{})
	go func() {
		log.Println("sleep 1s")
		time.Sleep(1 * time.Second)
		close(ch)
	}()
	CronUntil(func() {
		time.Sleep(100 * time.Millisecond)
		log.Println("test")
	}, 100*time.Millisecond, ch)
	log.Println("main exit")
}

type CusPod struct {
	ID   int
	Name string
}

func CronTry1() {
	podKillCh := make(chan *CusPod, 50)
	ch := make(chan struct{})

	go func() {
		i := 0
		for {
			time.Sleep(time.Second)
			podKillCh <- &CusPod{
				ID:   i,
				Name: strconv.Itoa(i),
			}
			i++
			if i == 5 {
				close(podKillCh)
				close(ch)
				return
			}
		}
	}()

	CronUntil(func() {
		for stu := range podKillCh {
			log.Printf("%+v\n", stu)
		}
	}, 1*time.Second, ch)

	log.Println("main exit")

}

func CronTry2() {

	ch := make(chan struct{})
	go func() {
		log.Println("sleep 1s")
		time.Sleep(10 * time.Second)
		close(ch)
	}()
	CronUntil(func() {
		time.Sleep(1000 * time.Millisecond)
		log.Println("test111")
	}, 2*time.Second, ch)
	log.Println("main exit")
}

func CronTry3() {

	NeverStopFunc(func() {
		time.Sleep(1000 * time.Millisecond)
		log.Println("test111")
	}, 2*time.Second)

}
