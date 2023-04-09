package producer_consumer_model

import (
	"fmt"
	"math/rand"
	"testing"
)

/*
	场景：多生产者+一生产聚合者+一消费者 加上退出通知。
	生产者：启动goroutine生产数据，并放入chan中，返回chan对象
	生产聚合者：用更大的chan来接收多个生产者的数据，并返回chan对象(select-for)
	消费者：主goroutine 从chan中拿出data处理
*/

// 建立多个Producer，并用chan传递data，再用一个函数 for-select merge起来。
// 加上退出通知机制。
func ProducerWithStop(done chan struct{}) chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("生产者1收到退出通知！")
				return
			case ch <-rand.Int():
				// 模拟处理业务逻辑。
			}

		}
	}()
	return ch
}

func ProducerWithStop2(done chan struct{}) chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("生产者2收到退出通知！")
				return
			case ch <-rand.Int():
				// 模拟处理业务逻辑。
			}

		}
	}()
	return ch

}

func MergeProducerWithStop(done chan struct{}) chan int {

	mergeCh := make(chan int, 50)
	go func() {
		for {
			select {
			case v1 := <-ProducerWithStop(done):
				mergeCh <- v1
			case v2 := <-ProducerWithStop2(done):
				mergeCh <- v2
			case <-done:
				fmt.Println("通知多个生产者们退出")
				fmt.Println("MergeProducer自己也退出！")
				return
			/*
				如果在select中执行send操作，则可能会永远被send阻塞。
				所以，在使用send的时候，应该也使用default语句块，保证send不会被阻塞。
				如果没有default，或者能确保select不阻塞的语句块，则迟早会被send阻塞。
			*/
			//default:
			//	fmt.Println("waiting")
			}
		}
	}()

	return mergeCh


}

func TestDemo3(t *testing.T) {
	doneC := make(chan struct{})
	ch := MergeProducerWithStop(doneC)
	for i := 0; i < 200; i++ {
		fmt.Println(<-ch)
	}

	doneC <- struct{}{}
	fmt.Println("主goroutine退出！")
}
