package timeout_mode

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestRetryTimeout(test *testing.T) {
	// 设置超时，超过十秒，整个链路会退出
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 永久重试
	//ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 重试并超时
	RetryTimeout(ctx, time.Second*2, errorRequest)
}

// errorRequest 模拟错误请求
func errorRequest(ctx context.Context) error {
	var body []byte

	a := http.DefaultClient
	a.Timeout = time.Second             // 设置超时时间
	res, err := a.Get("https://aaaaa/") // 随便写的url
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ = ioutil.ReadAll(res.Body)
	fmt.Println(body)
	return nil
}
