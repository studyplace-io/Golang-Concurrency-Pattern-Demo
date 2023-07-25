package worker_pool_mode

import (
	"log"
	"sync"
)

// LimitWaitGroup 限制并发数量的WaitGroup
type LimitWaitGroup struct {
	// WaitGroup对象
	Wg sync.WaitGroup
	// 设置的并发数
	size int
	// 使用chan的方式限制
	poolC chan struct{}
}

// LimitWaitGroupOption 使用配置模式来设置入参
type LimitWaitGroupOption func(nl *LimitWaitGroup)

const (
	defaultSize = 8
)

// NewDefaultLimitWaitGroup 兼容原来的WaitGroup
func NewDefaultLimitWaitGroup() *LimitWaitGroup {
	limitWaitGroup := &LimitWaitGroup{
		Wg:    sync.WaitGroup{},
		size:  0,
		poolC: make(chan struct{}, 0),
	}

	return limitWaitGroup
}

func NewLimitWaitGroup(opts ...LimitWaitGroupOption) *LimitWaitGroup {

	limitWaitGroup := &LimitWaitGroup{
		Wg:    sync.WaitGroup{},
		size:  defaultSize,
		poolC: make(chan struct{}, defaultSize),
	}

	for _, opt := range opts {
		opt(limitWaitGroup)
	}

	return limitWaitGroup
}

func WithSize(size int) LimitWaitGroupOption {
	return func(nl *LimitWaitGroup) {
		if size > 0 {
			nl.size = size
			nl.poolC = make(chan struct{}, size)
		} else {
			log.Fatal("size should > 0")
		}
	}
}

// BlockAdd 只是当size大于0时，
// 需要发一个struct{}给chan，相当于占用位置
func (nl *LimitWaitGroup) BlockAdd() {
	if nl.size > 0 {
		nl.poolC <- struct{}{}
	}
	nl.Wg.Add(1)
}

// Done 当size大于0时，需要给释放掉
func (nl *LimitWaitGroup) Done() {
	if nl.size > 0 {
		<-nl.poolC
	}
	nl.Wg.Done()
}

// Wait 方法
func (nl *LimitWaitGroup) Wait() {
	nl.Wg.Wait()
}

// PendingCount 返回当前状态LimitWaitGroup的数量
func (nl *LimitWaitGroup) PendingCount() int64 {
	return int64(len(nl.poolC))
}
