package producer_consumer_model

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

/*
	场景：多生产者+一生产聚合者+一消费者。
	生产者：启动goroutine生产数据，并放入chan中，返回chan对象
	生产聚合者：用更大的chan来接收多个生产者的数据，并返回chan对象(select-for)
	消费者：主goroutine 从chan中拿出data处理
*/

// Producer1 实现带有缓冲的生产者，返回一个chan，可从里面取出执行的结果
func Producer1() chan int {
	ch := make(chan int, 20)

	go func() {
		for {
			select {

			case ch <-rand.Int():
				// 执行业务逻辑

			}
		}
	}()

	return ch
}

// Producer2 实现带有缓冲的生产者，返回一个chan，可从里面取出执行的结果
func Producer2() chan int {
	ch := make(chan int, 20)

	go func() {
		for {
			select {
			default:
				// 执行业务逻辑
				ch <-rand.Int()
			}
		}
	}()

	return ch
}

// MergeProducer 聚合
// 建立多个Producer，并用chan传递data，再用一个函数 for-select merge起来。
func MergeProducer() chan int {

	ch := make(chan int, 30)

	go func() {
		for {
			select {
			case v1 := <-Producer1():
				ch <- v1
			case v2 := <-Producer2():
				ch <- v2
			}
		}
	}()

	return ch

}

func TestDemo2(t *testing.T) {
	// 聚合生产者
	ch := MergeProducer()

	// 消费者
	for i := 0; i < 200; i++ {
		fmt.Println(<-ch)
	}

	time.Sleep(time.Second)
}
