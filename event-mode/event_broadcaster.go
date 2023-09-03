package event_broadcaster

import (
	"fmt"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/event-mode/broadcaster"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/event-mode/event"
	"time"
)

// EventBroadcaster 事件广播器接口
type EventBroadcaster interface {
	// Event 发送事件给存储在广播器中的watcher对象
	Event(eType event.EventType , reason, message string)
	// EventBySource 增加事件源对象
	EventBySource(eType event.EventType , reason, message, source string)
	// Start 启动事件广播器
	Start()
	// Stop 停止事件广播器
	Stop()
}

// eventBroadcaster 实现事件广播器的对象
type eventBroadcasterImpl struct {
	// broadcaster 使用一个广播器实现
	*broadcaster.Broadcaster
	stopC    chan struct{}
}

// watcher queue
const queueLength = int64(1)

func NewEventBroadcaster() EventBroadcaster {
	return &eventBroadcasterImpl{
		broadcaster.NewBroadcaster(queueLength), make(chan struct{})}
}

func (eventBroadcaster *eventBroadcasterImpl) Stop() {
	eventBroadcaster.Shutdown()
	close(eventBroadcaster.stopC)
}

// Event 生成 event
func (eventBroadcaster *eventBroadcasterImpl) Event(eType event.EventType, reason, message string) {
	events := &event.Event{Type: eType, Reason: reason, Message: message, Timestamp: time.Now()}
	// 发送事件
	eventBroadcaster.Action(events)
}

// Event 生成 event
func (eventBroadcaster *eventBroadcasterImpl) EventBySource(eType event.EventType , reason, message, source string) {
	events := &event.Event{Type: eType, Reason: reason, Message: message, Source: source, Timestamp: time.Now()}
	// 发送事件
	eventBroadcaster.Action(events)
}

// Start 将日志打印
func (eventBroadcaster *eventBroadcasterImpl) Start() {
	// register a watcher
	watcher := eventBroadcaster.Watch()
	go func() {
		for watchEvent := range watcher.ResultChan() {
			fmt.Printf("%v\n", watchEvent)
		}
		<-eventBroadcaster.stopC
	}()
}
