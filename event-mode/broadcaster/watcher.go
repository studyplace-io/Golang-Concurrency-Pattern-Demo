package broadcaster

import (
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/event-mode/event"
)

// Interface 接口 watcher 实现
type Interface interface {
	Stop()
	ResultChan() <-chan event.Event
}

// broadcasterWatcher 观察者
type broadcasterWatcher struct {
	result  chan event.Event
	stopped chan struct{}
	id      int64
}

// 每个 watcher 通过该方法读取 channel 中广播的 events
func (b *broadcasterWatcher) ResultChan() <-chan event.Event {
	select {
	case <-b.stopped:
		return nil
	default:
		return b.result
	}

}

func (b *broadcasterWatcher) Stop() {
	b.stopped <- struct{}{}
}
