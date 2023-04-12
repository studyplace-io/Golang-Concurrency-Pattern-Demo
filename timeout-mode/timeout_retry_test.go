package timeout_mode

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// RetryTimeout 重试超时模式
func RetryTimeout(ctx context.Context, retryInterval time.Duration, execute func(ctx context.Context) error) {
	for {
		fmt.Println("execute func")
		if err := execute(ctx); err == nil {
			fmt.Println("work finished successfully")
			return
		}
		fmt.Println("execute if timeout has expired")
		if ctx.Err() != nil {
			fmt.Println("time expired 1 :", ctx.Err())
			return
		}
		fmt.Printf("wait %s before trying again\n", retryInterval)
		// 创建一个计时器
		t := time.NewTimer(retryInterval)
		select {
		case <-ctx.Done():
			fmt.Println("timed expired 2 :", ctx.Err())
			t.Stop()
			return
		// 定时执行！
		case <-t.C:
			fmt.Println("retry again")
		}
	}
}

func TestRetryTimeout(test *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 重试并超时
	RetryTimeout(ctx, time.Second*2, errorRequest)
}

// errorRequest 模拟错误请求
func errorRequest(ctx context.Context) error {
	var body []byte

	a := http.DefaultClient
	a.Timeout = time.Second  // 设置超时时间
	res, err  :=a.Get("https://aaaaa/") 	// 随便写的url
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ = ioutil.ReadAll(res.Body)
	fmt.Println(body)
	return nil
}
