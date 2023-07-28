package worker_pool_mode

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestLimitWaitGroup(test *testing.T) {

	urls := []string{
		"https://www.a.com/",
		"https://www.b.com",
		"https://www.c.com",
		"https://www.d.com/",
		"https://www.e.com",
		"https://www.f.com",
	}

	lp := NewLimitWaitGroup(WithSize(3))

	for _, url := range urls {
		lp.BlockAdd()
		url := url
		go func() {
			defer lp.Done()
			if url == "https://www.c.com" {
				time.Sleep(time.Duration(time.Second * 10))
			}
			if url == "https://www.a.com" {
				time.Sleep(time.Duration(time.Second * 10))
			}
			if url == "https://www.b.com" {
				time.Sleep(time.Duration(time.Second * 10))
			}

			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("error: %s: result: %v\n", url, err)
				return
			}
			defer resp.Body.Close()
		}()
	}

	// 模拟定时查看LimitWaitGroup状态
	go func() {
		n := time.NewTicker(time.Second)
		for {
			select {
			case <-n.C:
				fmt.Println("count goroutine: ", lp.PendingCount())
			default:
			}
		}
	}()

	lp.Wait()

	fmt.Println("Finished")
}
