package timeout_mode

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestContextCancel(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	stopC := make(chan struct{})
	for i := 0; i < 5; i++ {
		go doSomething(ctx, "child goroutine "+strconv.Itoa(i), stopC)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// 模拟过一定时间，定时关闭context
	go time.AfterFunc(time.Second*10, func() {
		defer wg.Done()
		fmt.Println("10 second, cancel func")
		cancel()
	})
	wg.Wait()

	fmt.Println("main goroutine done!")

}

func doSomething(ctx context.Context, name string, stopC chan struct{}) {
	i := 1
	for {

		time.Sleep(time.Second * 2)
		select {
		case <-ctx.Done():
			fmt.Printf("%s done!\n", name)
			fmt.Printf("%s 退出\n", name)
			//close(stopC)
			stopC <- struct{}{}
			return
		default:
			fmt.Printf("%s had worked %d seconds \n", name, i)

		}
		i++
	}

}

func TestContextTimeout(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stopC := make(chan struct{})

	for i := 0; i < 5; i++ {
		go doSomething(ctx, "child goroutine "+strconv.Itoa(i), stopC)
	}

	select {
	case <-stopC:
		fmt.Println("收到退出通知")
	}
	fmt.Println("main goroutine done!")

}

func TestContextTimeout2(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stopC := make(chan struct{})

	for i := 0; i < 5; i++ {
		go doSomething(ctx, "child goroutine "+strconv.Itoa(i), stopC)
	}

	select {
	case <-ctx.Done():
		fmt.Println("call successfully!!!")
		return
	case <-time.After(time.Duration(time.Second * 20)):
		fmt.Println("timeout!!!")
		return
	}

}
