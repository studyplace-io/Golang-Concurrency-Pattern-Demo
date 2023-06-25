package ring_buffer_mode

import "sync"

// ConcurrentQueue 支持并发的环形队列
type ConcurrentQueue interface {
	// Enqueue 并发入队方法
	Enqueue(data interface{}) error
	// Dequeue 并发出队方法
	Dequeue() (interface{}, error)
}

// CircularBuffer 环形队列，支持并发入队与出队操作，
// 当队列满或是队列空时，会阻塞当前goroutine
type circularBuffer struct {
	// 锁对象，Cond依赖此锁
	lock *sync.Mutex
	// 通知者
	notEmpty *sync.Cond
	notFull  *sync.Cond

	// taskQueue 队列
	taskQueue []any
	// capacity 队列大小
	capacity int
	// head 环形头
	head int
	// tail 环形尾
	tail int
}

func (s *circularBuffer) IsEmpty() bool {
	return s.head == s.tail
}

func (s *circularBuffer) IsFull() bool {
	return s.head == (s.tail+1)%s.capacity
}

func (s *circularBuffer) Enqueue(task any) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	// 检查是否队列满了
	for s.IsFull() {
		// 如果满了，需要阻塞goroutine
		s.notFull.Wait()
	}

	// 没有满的情况，修改成功入队且修改位置
	s.taskQueue[s.tail] = task
	s.tail = (s.tail + 1) % s.capacity

	// 处理如果队列为空，goroutine会阻塞，
	// 所以需要Signal通知一个goroutine
	s.notEmpty.Signal()

	return nil
}

func (s *circularBuffer) Dequeue() (any, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	// 检查是否队列空
	for s.IsEmpty() {
		// 如果空了，需要阻塞goroutine
		s.notEmpty.Wait()
	}

	data := s.taskQueue[s.head]
	s.head = (s.head + 1) % s.capacity

	// 处理如果队列满了，goroutine会阻塞，
	// 所以需要Signal通知一个goroutine
	s.notFull.Signal()

	return data, nil
}

func NewCircularBuffer(size int) ConcurrentQueue {
	w := &circularBuffer{
		taskQueue: make([]any, size),
		capacity:  size,
		lock:      &sync.Mutex{},
	}
	w.notEmpty = sync.NewCond(w.lock)
	w.notFull = sync.NewCond(w.lock)

	return w
}
