package future_mode

import (
	"io/ioutil"
	"net/http"
)

// RequestFuture 启动一个goroutine，请求http，并把结果放入chan中
func RequestFuture(url string) <-chan []byte {
	c := make(chan []byte, 1)
	go func() {
		var body []byte
		defer func() {
			c <- body
		}()

		res, err := http.Get(url)
		if err != nil {
			return
		}
		defer res.Body.Close()

		body, _ = ioutil.ReadAll(res.Body)
	}()

	return c
}

// RequestFutureV2 支持返回error结果
func RequestFutureV2(url string) (<-chan []byte, <-chan error) {
	c := make(chan []byte, 1)
	errC := make(chan error, 1)
	go func() {
		var body []byte
		defer func() {
			c <- body
		}()

		res, err := http.Get(url)
		if err != nil {
			errC <- err
			return
		}

		defer res.Body.Close()

		body, _ = ioutil.ReadAll(res.Body)
	}()

	return c, errC
}
