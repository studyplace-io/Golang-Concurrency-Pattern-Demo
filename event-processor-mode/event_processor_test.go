package event_processor_mode

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestEventProcessor(t *testing.T) {
	out := make(chan Event, 20)
	e := newEventProcessor(out)
	doneC := make(chan struct{})

	e.run()

	wg := sync.WaitGroup{}

	// 事件消费者，区分
	go func() {
		for {
			select {
			case <-doneC:
				return
			case o := <-out:
				// TODO 区分不同事件，可以做不同的处理
				switch o.Type {
				case Added:
					fmt.Println(o)
				case Modified:
					fmt.Println(o)
				case Deleted:
					fmt.Println(o)
				case Error:
					fmt.Println(o)
				}
			}
		}
	}()

	// 事件生产者
	for i := 0; i < 60; i++ {
		i := i
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := Event{
				Type: Added,
				Obj: i,
			}
			if i%3 == 0 {
				event.Type = Modified
				time.Sleep(time.Second)
			}
			// 推送事件
			e.push(event)
		}(i)
	}

	wg.Wait()
	e.stop()
	close(doneC)

}
