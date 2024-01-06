package event_broadcaster

import (
	"github.com/study-io/Golang-Concurrency-Pattern-Demo/event-mode/event"
	"testing"
	"time"
)

func TestEventBroadcaster(t *testing.T) {
	eventBroadcast := NewEventBroadcaster()

	eventBroadcast.Start()

	go func() {
		time.Sleep(time.Second * 3)
		eventBroadcast.Event(event.Warning, "test", "other-goroutine")
		time.Sleep(time.Second * 3)
		eventBroadcast.EventBySource(event.Normal, "test", "other-goroutine", "api-server")
		time.Sleep(time.Second * 3)
		eventBroadcast.Event(event.Normal, "test", "other-goroutine")
	}()

	time.Sleep(time.Second * 3)
	eventBroadcast.Event(event.Normal, "test", "main-goroutine")
	eventBroadcast.EventBySource(event.Normal, "test", "main-goroutine", "api-server")

	<-time.After(time.Second * 20)
	eventBroadcast.Stop()
}
