package ring_buffer_mode

import (
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {

	size := 20
	q := NewCircularBuffer(10)

	go func() {
		for i := 0; i < (size - 1); i++ {
			v, err := q.Dequeue()
			fmt.Println(v, err)
			time.Sleep(time.Second * 1)
		}
	}()

	for i := 0; i < (size - 1); i++ {
		go q.Enqueue(i)
	}

	time.Sleep(time.Second * 20)

}
