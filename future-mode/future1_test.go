package future_mode

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
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

func TestFuture(test *testing.T) {
	future := RequestFuture("https://api.github.com/users/octocat/orgs")
	// 中间可以实现自己的业务逻辑。。。。。
	fmt.Println("do something.....")

	// 当需的时候可以从chan出来
	body := <-future
	log.Printf("reponse length: %d", len(body))
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

func TestFutureWithError(test *testing.T) {

	res, errC := RequestFutureV2("https://api.github.com/users/octocat/orgs")

	// 中间可以实现自己的业务逻辑。。。。。
	fmt.Println("do something.....")

	select {
	case r := <-res:
		fmt.Println("res:", r)
	case e := <-errC:
		fmt.Println("err: ", e)
	}

}
