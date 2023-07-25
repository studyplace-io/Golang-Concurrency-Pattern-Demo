package timeout_mode

import (
	"context"
	"fmt"
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
