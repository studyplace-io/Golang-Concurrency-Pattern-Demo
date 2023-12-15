package broadcaster

import (
	"github.com/practice/Golang-Concurrency-Pattern-Demo/event-mode/event"
	"sync"
)

// Broadcaster 定义与实现
// 接收 events channel 的长度
const incomingQueueLength = 100

type Broadcaster struct {
	lock             sync.Mutex
	incoming         chan event.Event
	watchers         map[int64]*broadcasterWatcher
	watchersQueue    int64
	watchQueueLength int64
	distributing     sync.WaitGroup
}

func NewBroadcaster(queueLength int64) *Broadcaster {
	m := &Broadcaster{
		incoming:         make(chan event.Event, incomingQueueLength),
		watchers:         map[int64]*broadcasterWatcher{},
		watchQueueLength: queueLength,
	}
	m.distributing.Add(1)
	// 后台启动一个 goroutine 广播 events
	go m.loop()
	return m
}

// Watch 注册一个 watcher
func (m *Broadcaster) Watch() Interface {
	watcher := &broadcasterWatcher{
		result:  make(chan event.Event, incomingQueueLength),
		stopped: make(chan struct{}),
		id:      m.watchQueueLength,
	}
	m.watchers[m.watchersQueue] = watcher
	m.watchQueueLength++
	return watcher
}

// Action 接收所产生的 events
func (m *Broadcaster) Action(event *event.Event) {
	m.incoming <- *event
}

// 广播 events 到每个 watcher
func (m *Broadcaster) loop() {
	// 从 incoming channel 中读取所接收到的 events
	for event := range m.incoming {
		// 发送 events 到每一个 watcher
		for _, w := range m.watchers {
			select {
			case w.result <- event:
			case <-w.stopped:
			default:
			}
		}
	}
	// chan退出时的操作
	m.closeAll()
	m.distributing.Done()
}

// Shutdown 关闭广播器：1. 关闭chan 2. watigroup处理
func (m *Broadcaster) Shutdown() {
	close(m.incoming)
	m.distributing.Wait()
}

// closeAll 关闭广播器
func (m *Broadcaster) closeAll() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, w := range m.watchers {
		close(w.result)
	}
	m.watchers = map[int64]*broadcasterWatcher{}
}

func (m *Broadcaster) stopWatching(id int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	w, ok := m.watchers[id]
	if !ok {
		return
	}
	delete(m.watchers, id)
	close(w.result)
}
