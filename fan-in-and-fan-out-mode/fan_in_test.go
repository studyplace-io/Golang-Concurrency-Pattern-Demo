package fan_in_and_fan_out_mode

import (
	"fmt"
	"sync"
	"testing"
)

/*
 https://mp.weixin.qq.com/s?__biz=Mzg3MTA0NDQ1OQ==&mid=2247483680&idx=1&sn=de463ebbd088c0acf6c2f0b5f179f38d&scene=21#wechat_redirect
*/

// producer 负责生产数据，返回一个chan，把准备好的数据放入chan中。
func producer(num ...int) <-chan int {
	out := make(chan int)

	// 异步启goroutine准备数据，并放入chan
	go func() {
		defer close(out)
		for _, n := range num {
			out <- n
		}
	}()

	return out
}

// square 执行主要的业务逻辑，从准备好的chan中拿取数据，并执行业务逻辑，执行后放入chan中
func square(inputC <-chan int) <-chan int {

	out := make(chan int)
	// 异步启goroutine准备数据，并放入chan
	go func() {
		defer close(out)
		for n := range inputC {
			out <- n * n
		}
	}()
	return out

}

// merge FAN-IN 扇入模式
func merge(inputChans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// 从chan中把东西放入outC中
	collect := func(inputC <-chan int) {
		defer wg.Done()
		for n := range inputC {
			out <- n
		}
	}

	// FAN-IN，并发执行
	for _, c := range inputChans {
		wg.Add(1)
		go collect(c)
	}
	// FIXME 错误方式：直接等待是bug，死锁，因为merge写了out，main却没有读
	//wg.Wait()
	//close(out)

	// 正确方式
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func TestFanMode(t *testing.T) {
	in := producer(1, 3, 4, 5, 6)

	// 业务逻辑可能执行较长时间，可以使用并发
	c1 := square(in)
	c2 := square(in)
	c3 := square(in)
	// 把所有的chan放入一个总的chan中
	mergeC := merge(c1, c2, c3)

	for res := range mergeC {
		fmt.Println("res: ", res)
	}

	fmt.Println("fan mode finished.")
}
