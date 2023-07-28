package concurrent_queue_mode

import (
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	a := &queue{maxSize: 100, datas: make(chan interface{}, 100)}

	a.Put("aadddda")
	a.Put("aadddda")
	a.Put("aaddddadfasdkfj;af")
	a.Put("aadddsssssda")
	a.Put("adda")
	a.Put("aadddda")
	for i := 0; i < 10; i++ {
		aa, _ := a.Pop()
		fmt.Println(aa)
	}
}

func TestConcurrentQueue(t *testing.T) {

	c := NewConcurrentQueue(10)

	// 出队goroutine
	go func() {
		for {
			a, _ := c.Dequeue()
			fmt.Println("dequeue: ", a)
		}
	}()

	// 并发入队
	for i := 0; i < 100; i++ {
		go func(num int) {
			fmt.Println("num: ", num)
			c.Enqueue(num)
		}(i)
	}

	time.Sleep(time.Second * 3)
}
