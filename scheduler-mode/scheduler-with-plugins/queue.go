package scheduler_with_plugins

import _interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"

// Queue 调度队列
type Queue struct {
	// 所有加入队列的对象都放入此chan
	activeQ chan _interface.Pod
	// 当有调度错误时，放入backoffQ
	backoffQ  chan _interface.Pod

	// 调度时，从out chan中取对象
	out chan _interface.Pod
}

func NewScheduleQueue(capacity int) *Queue {
	return &Queue{
		activeQ: make(chan _interface.Pod, capacity),
		backoffQ:  make(chan _interface.Pod, capacity),
		out:       make(chan _interface.Pod, capacity*2),
	}
}

// Run 放入out chan的逻辑
// 一般情况：如果没有调度失败，就会把scheduleQ中的对象放入out中，
// 如果调度失败，就有概率从backoffQ中获取对象放入out中
func (q *Queue) Run(done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case t := <-q.backoffQ:
			q.out <- t
		default:
			select {
			case <-done:
				return
			case t := <-q.backoffQ:
				q.out <- t
			case t := <-q.activeQ:
				q.out <- t
			}
		}
	}
}

// Get first try backoffQ, and then scheduleQ. It blocks if both backoffQ and scheduleQ are empty
func (q *Queue) Get() <-chan _interface.Pod {
	return q.out
}

// Put 入队
func (q *Queue) Put(t _interface.Pod) {
	q.activeQ <- t
}

// Backoff 调度错误时，入队
func (q *Queue) Backoff(t _interface.Pod) {
	q.backoffQ <- t
}

func (q *Queue) Len() int {
	return len(q.activeQ) + len(q.backoffQ)
}
