package workqueue_mode

import (
	"fmt"
	"k8s.io/klog/v2"
	"testing"
	"time"
)

func TestSimpleQueue(t *testing.T) {

	q := newDelayingQueue(newQueue())
	// 加入回调方法
	q.SetCallback(CallbackFunc{
		AddFunc: func() {
			fmt.Println("something add in queue...")
		},
		GetFunc: func() {
			fmt.Println("something get from queue...")
		},
	})

	first := "test"
	first1 := "test1"
	first2 := "test2"
	first3 := "test3"
	first4 := "test4"
	// 延迟加入队列
	q.AddAfter(first, 50*time.Millisecond)
	q.AddAfter(first1, 50*time.Millisecond)
	q.AddAfter(first2, 50*time.Millisecond)
	q.AddAfter(first4, 80*time.Millisecond)
	q.AddAfter(first3, 100*time.Millisecond)

	// 确认队列中目前没有元素
	if q.Len() != 0 {
		t.Errorf("should not have added")
	}

	time.Sleep(70 * time.Millisecond)
	klog.Infof("pass 70 millisecond, queue len: %v", q.Len())
	time.Sleep(100 * time.Millisecond)
	klog.Infof("pass 100 millisecond, queue len: %v", q.Len())
	item, _ := q.Get()
	q.Done(item)
	klog.Infof("get item:  %v", item)
	klog.Infof("after get ont item, queue len: %v", q.Len())

	q.AddAfter("test-add-in-8-seconds", 8*time.Second)
	q.AddAfter("test2-add-in-8-seconds", 8*time.Second)
	q.AddAfter("test3-add-in-8-seconds;fj", 8*time.Second)

	// 经过一次 heartbeat
	time.Sleep(12 * time.Second)
	item1, _ := q.Get()
	q.Done(item1)
	klog.Infof("get item:  %v", item1)
	klog.Infof("final, queue len: %v", q.Len())
	q.Close()
}
