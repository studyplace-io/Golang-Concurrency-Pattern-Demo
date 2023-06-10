package concurrent_queue_mode

// Queue 队列
type Queue interface {
	// Put 入队
	Put(data interface{}) error
	// Pop 出队
	Pop() (interface{}, error)
	IsEmpty() bool
	IsFull() bool
	SetLen(uint32)
	SetCap(uint32)
	InitQueue()
	Len() uint32
	Cap() uint32
}

// ConcurrentQueue 支持并发队列
type ConcurrentQueue interface {
	// Enqueue 并发入队方法
	Enqueue(data interface{}) error
	// Dequeue 并发出队方法
	Dequeue() (interface{}, error)
}
