package timeout_mode

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWaitGroupWithoutTimeout(test *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()
			fmt.Println(index)
			// 模拟其中一个任务执行的非常久，所有goroutine都在等
			if index == 2 {
				time.Sleep(time.Second * 200)
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("所有任务执行完毕。。。。")

}

func TestWaitGroupWithChan(test *testing.T) {
	var wg sync.WaitGroup
	doneC := make(chan struct{})

	go func() {
		wg.Wait()
		doneC <- struct{}{}
	}()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fmt.Println(index)
			// 模拟某个goroutine 执行非常久的时间
			if index == 2 {
				time.Sleep(time.Second * 100)
			}
		}(i)
	}

	timeout := time.Duration(10) * time.Second
	fmt.Printf("Wait for waitgroup (up to %s)\n", timeout)

	select {
	case <-doneC:
		fmt.Printf("Wait group finished\n")
	case <-time.After(timeout):
		fmt.Printf("Timed out waiting for wait group\n")
	}

}

func TestWaitGroupWithContext(test *testing.T) {
	ww := NewWaitGroupWithTimeout(time.Second * 5)
	c := make(chan struct{})
	for i := 0; i < 10; i++ {
		ww.Add(1)
		numC := make(chan int, 10)
		i := i
		numC <- i
		go func(numC chan int, close <-chan struct{}) {
			defer ww.Done()
			for {
				select {
				case <-close:
					fmt.Println("超时退出")
					return
				case m := <-numC:
					fmt.Println(m)
					if m == 2 {
						time.Sleep(time.Second * 2)
					}
				}

			}

		}(numC, c)
	}

	if ww.WaitTimeout() {
		close(c)
		fmt.Println("timeout exit")
	}

}

// WaitGroupWithTimeout 增加超时功能的WaitGroup
type WaitGroupWithTimeout struct {
	sync.WaitGroup
	Timeout time.Duration // 超时时间
}

func NewWaitGroupWithTimeout(timeout time.Duration) *WaitGroupWithTimeout {
	w := &WaitGroupWithTimeout{
		sync.WaitGroup{},
		timeout,
	}
	return w
}

// WaitTimeout 判断Wait() 是否超时
func (wg *WaitGroupWithTimeout) WaitTimeout() bool {

	ch := make(chan bool, 1)

	go time.AfterFunc(wg.Timeout, func() {
		ch <- true
	})

	go func() {
		wg.Wait()
		ch <- false
	}()

	return <-ch
}