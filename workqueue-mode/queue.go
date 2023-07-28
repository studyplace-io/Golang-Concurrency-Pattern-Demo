package workqueue_mode

import (
	"sync"
)

// Interface 队列接口对象
type Interface interface {
	// Add 加入工作队列
	Add(item interface{})
	// Get 从工作队列中取出对象
	Get() (item interface{}, shutdown bool)
	// Done 记录该对象已经被处理
	Done(item interface{})
	SetCallback(handler CallbackHandler)
	Len() int
	ShutDown()
	IsShutDown() bool
}

// obj 放入队列的对象
type obj interface{}

type queue struct {
	// queue 存储对象的对列
	queue []obj
	// dirty set 用来去重，避免同个对象被重复加入多次
	dirty set
	// 标记是否正在被处理
	processing set
	cond       *sync.Cond
	// 队列是否关闭状态
	close bool

	// CallbackHandler 当缓存出现修改时，可执行的回调方法
	callbacks CallbackHandler
}

// CallbackHandler 回调接口，可提供用户实现相应方法
type CallbackHandler interface {
	OnAdd()
	OnGet()
}

// CallbackFunc 回调方法
type CallbackFunc struct {
	// OnAdd 加入队列，可执行的回调
	AddFunc func()
	// OnGet 获取队列，可执行的回调
	GetFunc func()
}

func (c CallbackFunc) OnAdd() {
	if c.AddFunc != nil {
		c.AddFunc()
	}
}

func (c CallbackFunc) OnGet() {
	if c.GetFunc != nil {
		c.GetFunc()
	}
}

func newQueue() *queue {
	t := &queue{
		queue:      make([]obj, 0),
		dirty:      set{},
		processing: set{},
		cond:       sync.NewCond(&sync.Mutex{}),
	}
	return t
}

// set 元组用来去重
type set map[obj]struct{}

func (s set) has(item obj) bool {
	_, exists := s[item]
	return exists
}

func (s set) insert(item obj) {
	s[item] = struct{}{}
}

func (s set) delete(item obj) {
	delete(s, item)
}

// Add 加入工作队列
func (q *queue) Add(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	// 如果关闭，就直接返回
	if q.close {
		return
	}
	// 如果已经加入，不要重复加入
	if q.dirty.has(item) {
		return
	}

	// 加入dirty
	q.dirty.insert(item)

	// 如果已经在处理，代表并发情况，不要直接加入queue中，
	//等待正在处理的那个goroutine处理完后，会加入。
	if q.processing.has(item) {
		return
	}

	if q.callbacks != nil {
		q.callbacks.OnAdd()
	}
	// 加入queue中
	q.queue = append(q.queue, item)
	q.cond.Signal()
}

func (q *queue) Len() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return len(q.queue)
}

// Get 从工作队列获取
func (q *queue) Get() (item interface{}, shutdown bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	// 如果队列中没有，等待
	for len(q.queue) == 0 && !q.close {
		q.cond.Wait()
	}
	if len(q.queue) == 0 && q.close {
		// We must be shutting down.
		return nil, true
	}

	if q.callbacks != nil {
		q.callbacks.OnGet()
	}

	// 变化切片
	item, q.queue = q.queue[0], q.queue[1:]

	// 加入正在处理的标示
	q.processing.insert(item)
	// 从标示中取出
	q.dirty.delete(item)

	return item, false
}

// Done 当处理完后，需要调用
func (q *queue) Done(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	// 从正在处理标示中删除
	q.processing.delete(item)
	// 如果标示中有，代表是并发情况，需要加入到queue中
	if q.dirty.has(item) {
		q.queue = append(q.queue, item)
		q.cond.Signal()
	}
}

func (q *queue) SetCallback(handler CallbackHandler) {
	q.callbacks = handler
}

func (q *queue) ShutDown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.close = true
	q.cond.Broadcast()
}

func (q *queue) IsShutDown() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.close
}
