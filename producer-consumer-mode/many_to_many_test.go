package producer_consumer_mode

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

/*
	多个生产者+多个消费者的模式
	不要close 缓存的dataCh，使用stopC从消费者端通知。

	注：不要全部都用waitgroup来做，只要消费者用即可。
	因为生产者有close兜底了，所以只要确认消费者的退出就好。
*/

// 常数
const (
	Max        = 100000
	NumSenders = 10
)

var wgReceivers3 sync.WaitGroup

func ProducerConsumerManyToMany() {

	rand.Seed(time.Now().UnixNano())

	// 数据缓存
	dataC := make(chan int, 100)
	// 通知退出
	stopC := make(chan struct{})

	// 多个生产者
	for i := 0; i < NumSenders; i++ {
		go ProducerMany(dataC, stopC, i)
	}

	// 消费者，多个消费者需要使用waitgroup
	wgReceivers3.Add(1)
	go Consumer(dataC, &wgReceivers3, stopC)

	wgReceivers3.Wait()

}

func TestProducerConsumerManyToMany(t *testing.T) {

	ProducerConsumerManyToMany()

}

// ProducerMany 实现带有缓冲的生产者，返回一个chan，可从里面取出执行的结果
func ProducerMany(ch chan int, stopC chan struct{}, num int) chan int {

	fmt.Printf("第%v个生产者\n", num)

	go func(stopC chan struct{}) {
		for {
			select {
			case <-stopC:
				fmt.Printf("生产者%v退出\n", num)
				time.Sleep(time.Second)
				return
			// TODO 这里可以处理生产的业务逻辑
			// 。。。。。
			case ch <- rand.Intn(Max):

			}
		}
	}(stopC)

	return ch
}

// Consumer 消费者，做两件事情：
// 1. 取出chan中的数据，并消费。
// 2. 提供退出逻辑，通知生产者们。
func Consumer(dataC chan int, wgReceivers *sync.WaitGroup, stopC chan struct{}) {

	defer wgReceivers.Done()
	// 遍历取出来消费
	for value := range dataC {
		// TODO 这里可以处理消费的业务逻辑
		// 。。。。。

		// 提供一个退出的条件，由消费者来控制退出！
		if value == Max-1 {
			fmt.Println("send stop signal to senders.")
			close(stopC)
			return
		}
		fmt.Println("value:", value)
	}

}
