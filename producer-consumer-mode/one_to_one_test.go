package producer_consumer_mode

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

/*
	场景：一生产者+一消费者。
	生产者：启动goroutine生产数据，并放入chan中，返回chan对象
	消费者：主goroutine 从chan中拿出data处理
*/

// Producer 实现带有缓冲的生产者，返回一个chan，可从里面取出执行的结果
func Producer(stopC chan struct{}) chan int {
	ch := make(chan int, 20)

	go func() {
		for {
			select {
			case <-stopC:
				log.Println("收到退出消息，生产者退出")
				return
			case ch <- rand.Int():

			}
		}
	}()

	return ch
}

func TestProducer(t *testing.T) {
	stopC := make(chan struct{})
	ch := Producer(stopC)
	// 模拟消费者
	for i := 0; i < 100; i++ {
		fmt.Println(<-ch)
	}
	// 通知关闭信号，关闭传出的chan
	stopC <- struct{}{}
	close(ch)
	// 测试过一段时间后还继续从chan中取出数据
	time.Sleep(time.Second * 10)
	for i := 0; i < 20; i++ {
		fmt.Println(<-ch)
	}
}
