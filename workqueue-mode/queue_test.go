package workqueue_mode

import (
	"k8s.io/klog/v2"
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {

	q := newQueue()
	// Start producers
	const producers = 50
	producerWG := sync.WaitGroup{}
	producerWG.Add(producers)
	for i := 0; i < producers; i++ {
		go func(i int) {
			defer producerWG.Done()
			for j := 0; j < 50; j++ {
				q.Add(i)
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	// Start consumers
	const consumers = 10
	consumerWG := sync.WaitGroup{}
	consumerWG.Add(consumers)
	for i := 0; i < consumers; i++ {
		go func(i int) {
			defer consumerWG.Done()
			for {
				item, quit := q.Get()
				if item == "added after shutdown!" {
					t.Errorf("Got an item added after shutdown.")
				}
				if quit {
					return
				}
				t.Logf("Worker %v: begin processing %v", i, item)
				time.Sleep(3 * time.Millisecond)
				t.Logf("Worker %v: done processing %v", i, item)
				q.Done(item)
			}
		}(i)
	}

	producerWG.Wait()
	q.ShutDown()
	q.Add("added after shutdown!")

	consumerWG.Wait()
}

func TestAddWhileProcessing(t *testing.T) {
	q := newQueue()

	// Start producers
	const producers = 50
	producerWG := sync.WaitGroup{}
	producerWG.Add(producers)
	for i := 0; i < producers; i++ {
		go func(i int) {
			defer producerWG.Done()
			q.Add(i)
		}(i)
	}

	// Start consumers
	const consumers = 10
	consumerWG := sync.WaitGroup{}
	consumerWG.Add(consumers)
	for i := 0; i < consumers; i++ {
		go func(i int) {
			defer consumerWG.Done()
			// Every worker will re-add every item up to two times.
			// This tests the dirty-while-processing case.
			counters := map[interface{}]int{}
			for {
				item, quit := q.Get()
				if quit {
					return
				}
				klog.Info(item)
				counters[item]++
				if counters[item] < 2 {
					q.Add(item)
				}
				q.Done(item)
			}
		}(i)
	}

	producerWG.Wait()
	q.ShutDown()
	consumerWG.Wait()
}
