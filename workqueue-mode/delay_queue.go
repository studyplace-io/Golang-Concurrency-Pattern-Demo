package workqueue

import (
	"container/heap"
	"k8s.io/klog/v2"
	"time"
)

// DelayQueue 延迟队列
type DelayQueue interface {
	// Interface 队列接口对象
	Interface
	// AddAfter 延迟duration时间入队
	AddAfter(item obj, duration time.Duration)
	Close()
}

// delayingQueue 延时队列实现对象
type delayingQueue struct {
	// 继承 Queue Interface 队列基本功能
	Interface
	// 退出通道
	stopCh chan struct{}

	// 堆对象，使用切片实现
	priorityQueue *PriorityQueue
	// waitingForAddCh 新的定时元素会加入chan, 等待 loop 处理.
	heapObjAddCh chan *heapObj
	// knownEntries 记录此元素是否被堆记录过
	knownEntries map[obj]*heapObj

	// 周期性检测队列是否有对象到期
	heartbeat *time.Ticker
}

// 心跳的时长
const maxWait = 10 * time.Second

func NewDelayingQueue(q Interface) DelayQueue {
	ret := &delayingQueue{
		Interface:     q,
		heartbeat:     time.NewTicker(maxWait),
		stopCh:        make(chan struct{}),
		heapObjAddCh:  make(chan *heapObj, 1000),
		knownEntries:  map[obj]*heapObj{},
		priorityQueue: &PriorityQueue{},
	}

	// 执行堆元素处理器
	go ret.heapHandlerLoop()
	return ret
}

// heapHandlerLoop 启动堆处理器，实现计时插入队列功能
// 1. 初始化对堆对象
// 2. 不断循环操作
func (q *delayingQueue) heapHandlerLoop() {

	// 初始化 min heap 小顶堆
	heap.Init(q.priorityQueue)

	for {
		// 如果 queue 已经被关闭, 则退出该 loop 协程.
		if q.Interface.IsShutDown() {
			return
		}

		now := time.Now()

		// 如果堆中有元素，需要执行的操作
		for q.priorityQueue.Len() > 0 {
			// 获取堆顶的元素
			entry := q.priorityQueue.Peek().(*heapObj)
			// 如果大于当前时间, 则没有到期, 则跳出.
			if entry.readyAt.After(now) {
				break
			}

			// 如果小于当前时间, 则 pop 出元素, 然后加入 queue 队里中.
			entry = heap.Pop(q.priorityQueue).(*heapObj)
			q.Add(entry.data)
			delete(q.knownEntries, entry.data)
		}

		// 如果小顶堆为空, 则使用 never 做无限时长定时器
		nextReadyAt := make(<-chan time.Time)

		var nextReadyAtTimer *time.Timer

		// 如果 堆不为空, 设置最近元素的时间为定时器的时间.
		if q.priorityQueue.Len() > 0 {
			// 把之前的Timer对象取消
			if nextReadyAtTimer != nil {
				nextReadyAtTimer.Stop()
			}

			// 从堆顶获取最近的元素
			entry := q.priorityQueue.Peek().(*heapObj)

			// 实例化 timer 定时器，把最接近下次的对象设置一个定时器，由下面select触发
			nextReadyAtTimer = time.NewTimer(entry.readyAt.Sub(now))
			nextReadyAt = nextReadyAtTimer.C
		}

		select {
		case <-q.stopCh:
			return
		// 触发 10s 心跳超时后, 重新进行选择最近的定时任务.
		case <-q.heartbeat.C:
			klog.Info("delay queue heartbeat...")
		// 最近元素的定时器已到期, 进行下次循环. 期间会处理该到期任务.
		case <-nextReadyAt:
			klog.Info("delay queue next obj is about to be executed ...")
		// 收到新添加的定时器
		// 如果新对象还未到期, 则把定时对象放到 heap 定时堆里.
		case waitEntry := <-q.heapObjAddCh:
			if waitEntry.readyAt.After(time.Now()) {
				q.insert(waitEntry)
			} else {
				// 如果该定时任务已到期, 则调用继承的 queue 的 add 方法.把元素添加到队列中.
				q.Add(waitEntry.data)
			}
		}
	}
}

// AddAfter 调用方使用 AddAfter 延迟加入队列
func (q *delayingQueue) AddAfter(item obj, duration time.Duration) {
	// 如果队列状态为关闭, 直接退出
	if q.IsShutDown() {
		return
	}

	// 如果时间不为正值, 直接入队列，不延迟入队
	if duration <= 0 {
		q.Add(item)
		return
	}

	select {
	// 等待退出
	case <-q.stopCh:
	// 创建一个定时对象, 然后推到 waitingForAddCh 管道中, 等待 waitingLoop 协程处理.
	case q.heapObjAddCh <- &heapObj{data: item, readyAt: time.Now().Add(duration)}:
	}
}

// Close 关闭延迟队列
func (q *delayingQueue) Close() {
	q.ShutDown()
	q.stopCh <- struct{}{}
}

// insert 把对象加入堆中, 如果已经存在，则更新readyAt时间
func (q *delayingQueue) insert(entry *heapObj) {

	existing, exists := q.knownEntries[entry.data]
	if exists {
		if existing.readyAt.After(entry.readyAt) {
			existing.readyAt = entry.readyAt
			// 调用堆接口，修改位置
			heap.Fix(q.priorityQueue, existing.index)
		}

		return
	}
	// 调用堆接口，push进入
	heap.Push(q.priorityQueue, entry)
	q.knownEntries[entry.data] = entry
}

type heapObj struct {
	// 存放元素
	data obj
	// 时间点
	readyAt time.Time
	// 同一个时间点下对比递增的索引
	index int
}

// PriorityQueue 小顶堆实现定时器
type PriorityQueue []*heapObj

func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].readyAt.Before(pq[j].readyAt)
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*heapObj)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	item.index = -1
	*pq = (*pq)[0:(n - 1)]
	return item
}

func (pq PriorityQueue) Peek() interface{} {
	return pq[0]
}
