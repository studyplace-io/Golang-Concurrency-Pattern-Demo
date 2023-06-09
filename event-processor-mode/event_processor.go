package event_processor_mode

import (
	"sync"
	"time"
)

// EventType 定义不同事件类型
type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Error    EventType = "ERROR"
)

// Event 事件
type Event struct {
	// 事件类型
	Type EventType
	// 需要传递的对象
	Obj interface{}
}

func newEventProcessor(out chan<- Event) *eventProcessor {
	return &eventProcessor{
		out:  out,
		cond: sync.NewCond(&sync.Mutex{}),
		done: make(chan struct{}),
	}
}

type EventProcessor interface {
	run()
	takeBatch() []Event
	writeBatch(events []Event)
	push(event Event)
	stop()
}

// eventProcessor 事件通知器
type eventProcessor struct {
	// 下发事件通知的chan
	out chan<- Event
	// 并发通知
	cond *sync.Cond
	// buff 处理事件是批量执行的，用来存储一批事件
	buff []Event

	done chan struct{}
}

// run 执行
func (e *eventProcessor) run() {
	go func() {
		for {
			// 取出一定批量的消息
			batch := e.takeBatch()
			// 执行写入chan
			e.writeBatch(batch)
			if e.stopped() {
				return
			}
		}
	}()
}

// takeBatch 批量执行
func (e *eventProcessor) takeBatch() []Event {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()

	// 如果没有，先阻塞
	for len(e.buff) == 0 && !e.stopped() {
		e.cond.Wait()
	}

	// 把值赋给batch后，重新开始新的buff
	batch := e.buff
	e.buff = nil
	return batch
}

// writeBatch 把一批消息写入chan
func (e *eventProcessor) writeBatch(events []Event) {
	for _, event := range events {
		select {
		case e.out <- event:
		case <-e.done:
			return
		}
	}
}

// push 放入消息
func (e *eventProcessor) push(event Event) {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()
	defer e.cond.Signal()
	e.buff = append(e.buff, event)
}

// stopped 停止
func (e *eventProcessor) stopped() bool {
	select {
	case <-e.done:
		return true
	default:
		return false
	}
}

func (e *eventProcessor) stop() {
	// 停止前，先检查buff中是否还有数据。
	for len(e.buff) != 0 {
		time.Sleep(time.Millisecond)
	}
	close(e.done)
	e.cond.Signal()
}
