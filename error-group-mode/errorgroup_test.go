package error_group_mode

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
	"testing"
)


/*
    errgroup的使用实践。

	https://mp.weixin.qq.com/s/qqva2Tj2qijWv_YCplWZ8A
	1. 继承了 WaitGroup 的功能
	2. 错误传播：能够返回任务组中发生的第一个错误，但有且仅能返回该错误
	3. context 信号传播：如果子任务 goroutine 中有循环逻辑，则可以添加 ctx.Done 逻辑，此时通过 context 的取消信号，提前结束子任务执行。
*/


func TestErrorPractice1(t *testing.T) {
	//TryUseWaitGroup()
	//TryUseErrGroup()
	ErrGroupUseContext()
}

// TryUseWaitGroup 使用waitGroup实现
// 模拟请求url，没有使用chan error来接住所有goroutine的error。
func TryUseWaitGroup() {

	var urls = []string{
		"http://www.golang.org/",
		"http://www.baidu.com/",
		"http://www.noexist11111111.com/",
	}
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			resp.Body.Close()
		}(url)
	}

	wg.Wait()
	fmt.Println("goroutine执行完毕，结束阻塞")
}

// TryUseErrGroup 使用errgroup.Group来捞出有错的goroutine
func TryUseErrGroup() {

	var urls = []string{
		"http://www.golang.org/",
		"http://www.baidu.com/",
		"http://www.noexist11111111.com/",
		"http://www.noexist11111111.com/",
		"http://www.noexist11111111.com/",
		"http://www.noexist11111111.com/",
	}

	g := new(errgroup.Group)

	for _, url := range urls {
		url := url
		g.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			fmt.Printf("get [%s] success: [%d] \n", url, resp.StatusCode)
			return resp.Body.Close()

		})
	}

	// 注意：这里的wait是 "当所有的goroutine中，只要有一个报错，就会返回且退出"
	if err := g.Wait(); err != nil {
		// 只会返回"第一个error"，不管后续还有没有error，都不再执行
		fmt.Println(err)
	} else {
		fmt.Println("All success!")
	}
	fmt.Println("主goroutine阻塞结束，进程退出")
}

func ErrGroupUseContext() {
	// 创建一个context上下文
	ctx := context.Background()
	// 创建 errgroup的上下文
	g, ctx := errgroup.WithContext(ctx)

	// 通道 存放data
	dataChan := make(chan int, 20)

	// 启一个生产者
	g.Go(func() error {
		// 单一个生产者，用完chan记得可以关闭
		defer close(dataChan)

		// 不断增加
		for i := 1; ;i++ {
			// 到了特定条件，oupput错误。这时候 errgroup.WithContext
			if i % 2 == 0 {
				return fmt.Errorf("data %v is wrong", i)
			}
			// 这里可以执行业务逻辑。
			fmt.Println(fmt.Sprintf("sending %d", i))

			// 把对象放进chan
			dataChan <- i

		}
	})

	// 启三个消费者
	for i := 0; i < 3; i++ {
		// 启Goroutine
		g.Go(func() error {
			for j := 1; ; j++ {
				// 监听不同chan
				select {
				// 当有错误，把错误返回！ 直接退出
				case <- ctx.Done():
					return ctx.Err()
				// 正常情况，接收到data，执行某些业务逻辑
				case num := <-dataChan:
					fmt.Println(fmt.Sprintf("receiving %d", num))
				}
			}
		})
	}

	// 这里会阻塞，如果有err，会打印err
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("main goroutine done!")

}
