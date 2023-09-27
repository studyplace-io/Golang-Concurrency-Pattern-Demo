package pubsub

import (
	"log"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	b := NewBroker(8)
	go b.Run()

	consumer := func(name string) ConsumeFunc {
		return func(topic Topic) {
			log.Printf("[%v] receive topic: %+v", name, topic)
		}
	}

	s1 := NewSubscriber("s1", b)
	s1.Subscribe("t1", consumer("s1"))
	s1.Subscribe("t2", consumer("s1"))

	s2 := NewSubscriber("s2", b)
	s2.Subscribe("t2", consumer("s2"))

	s3 := NewSubscriber("s3", b)
	s3.Subscribe("t3", consumer("s3"))

	p1 := NewPublisher("p1", b)
	p2 := NewPublisher("p2", b)

	p1.Publish(Topic{TopicID: "t1", Message: "hello xxx"})
	p1.Publish(Topic{TopicID: "t2", Message: "hello xxx2"})
	p1.Publish(Topic{TopicID: "t3", Message: "hello xxx3"})

	p2.Publish(Topic{TopicID: "t2", Message: "p2 hello xxx"})

	s1.Unsubscribe("t2")
	p1.Publish(Topic{TopicID: "t2", Message: "hello again xxx"})

	time.Sleep(2 * time.Second)
}
