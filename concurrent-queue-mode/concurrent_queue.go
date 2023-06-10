package concurrent_queue_mode

import (
	"sync"
)

// concurrentQueue 并发队列
type concurrentQueue struct {
	// mutex lock
	lock *sync.Mutex

	// empty and full locks
	notEmpty *sync.Cond
	notFull  *sync.Cond

	// queue storage backend
	backend Queue
}

// Enqueue 并发入队
func (c *concurrentQueue) Enqueue(data interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 如果队列满，所有goroutine阻塞
	for c.backend.IsFull() {
		//wait for empty
		c.notFull.Wait()
	}

	// 放入元素
	err := c.backend.Put(data)

	// 通知
	c.notEmpty.Signal()

	return err
}

// Dequeue 并发出队
func (c *concurrentQueue) Dequeue() (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 如果队列空，所有goroutine阻塞
	for c.backend.IsEmpty() {
		c.notEmpty.Wait()
	}

	data, err := c.backend.Pop()

	// 通知
	c.notFull.Signal()

	return data, err
}

func (c *concurrentQueue) getSize() uint32 {
	c.lock.Lock()
	defer c.lock.Unlock()
	size := c.backend.Len()

	return size
}

// NewConcurrentQueue Creates a new queue
func NewConcurrentQueue(maxSize uint32) ConcurrentQueue {
	cQueue := concurrentQueue{}

	//init mutexes
	cQueue.lock = &sync.Mutex{}
	cQueue.notFull = sync.NewCond(cQueue.lock)
	cQueue.notEmpty = sync.NewCond(cQueue.lock)

	//init backend
	cQueue.backend = &queue{}
	cQueue.backend.SetLen(0)
	cQueue.backend.SetCap(maxSize)
	cQueue.backend.InitQueue()

	return &cQueue
}
