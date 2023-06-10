package concurrent_queue_mode

import "errors"

// Queue 队列对象，底层使用chan实现
type queue struct {
	// datas 存储chan
	datas chan interface{}
	// size 目前存储容量
	size uint32
	// maxSize 队列最大容量
	maxSize uint32
}

// put 放入队列
func (queue *queue) Put(data interface{}) error {
	// 1. 检查是否超过最大容量
	if queue.size >= queue.maxSize {
		return errors.New("Queue full")
	}
	// 入队
	queue.datas <- data
	queue.size++
	return nil
}

// pop 弹出队列
func (queue *queue) Pop() (interface{}, error) {
	// 如果容量为0
	if queue.size == 0 {
		return nil, errors.New("Queue empty")
	}
	queue.size--
	return <-queue.datas, nil
}


func (queue *queue) IsEmpty() bool {
	return queue.size == 0
}


func (queue *queue) IsFull() bool {
	return queue.size >= queue.maxSize
}

func (queue *queue) Len() uint32 {
	return queue.size
}

func (queue *queue) Cap() uint32 {
	return queue.maxSize
}

func (queue *queue) SetLen(size uint32) {
	queue.size = size
}

func (queue *queue) SetCap(maxSize uint32) {
	queue.maxSize = maxSize
}

func (queue *queue) InitQueue() {
	queue.datas = make(chan interface{}, queue.maxSize)
}