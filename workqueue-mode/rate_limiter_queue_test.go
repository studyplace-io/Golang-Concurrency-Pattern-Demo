package workqueue

import (
	"fmt"
	"testing"
	"time"
)

func TestRateLimitingQueue(t *testing.T) {
	opts := RateLimitingQueueOption{rate: 1, buckets: 2}
	q := NewRateLimitingQueue(opts)

	q.AddRateLimited("one")
	q.AddRateLimited("one1")
	q.AddRateLimited("one2")
	q.AddRateLimited("two")
	q.AddRateLimited("two1")

	for i := 0; i < 1000; i++ {
		item, _ := q.Get()
		q.Done(item)
		fmt.Println(item)
		time.Sleep(time.Second)
		if q.Len() == 0 {
			break
		}
	}

	q.Close()
}
