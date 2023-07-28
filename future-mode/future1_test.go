package future_mode

import (
	"fmt"
	"log"
	"testing"
)

func TestFuture(test *testing.T) {
	future := RequestFuture("https://api.github.com/users/octocat/orgs")
	// 中间可以实现自己的业务逻辑。。。。。
	fmt.Println("do something.....")

	// 当需的时候可以从chan出来
	body := <-future
	log.Printf("reponse length: %d", len(body))
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
