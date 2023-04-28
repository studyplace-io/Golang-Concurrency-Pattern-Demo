package scheduler_with_plugins

import (
	"fmt"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/demo"
	_interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/nodes"
	"k8s.io/klog/v2"
	"sync"
	"testing"
	"time"
)

func TestSchedulerWithPlugin(test *testing.T) {

	schedulerName := fmt.Sprintf("%s-scheduler", "test")

	// 调度器
	scheduler := Scheduler{
		name:      schedulerName,
		options:   &defaultOptions,
		queue:     NewScheduleQueue(defaultOptions.queueCapacity),
		nodeInfos: nodes.NewNodeInfos(),
		workers:   defaultOptions.numWorker,
		plugins:   make([]_interface.Plugin, 0),
		wg:        sync.WaitGroup{},
		stopC:     make(chan struct{}),
		logger:    klog.LoggerWithName(klog.Background(), schedulerName),
	}

	// 加入模拟node
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node1"))
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node2"))
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node3"))

	// 启动调度器
	scheduler.wg.Add(1)
	go scheduler.run()

	go func() {
		// 需要等待并发
		scheduler.wg.Wait()
	}()

	defer func() {
		// 通知退出
		scheduler.Stop()
		time.Sleep(time.Second * 2)
	}()

	// 加入pod对象入队

	res := make([]*demo.MockPod, 0)

	t1 := &demo.MockPod{
		Name:          "pod1",
		PodRecordNode: make([]_interface.PodRecordNode, 0),
	}
	scheduler.AddPod(t1)
	res = append(res, t1)

	t2 := &demo.MockPod{
		Name:          "pod1",
		PodRecordNode: make([]_interface.PodRecordNode, 0),
	}
	scheduler.AddPod(t2)
	res = append(res, t2)

	t3 := &demo.MockPod{
		Name:          "pod1",
		PodRecordNode: make([]_interface.PodRecordNode, 0),
	}
	scheduler.AddPod(t3)
	res = append(res, t3)

	// 加入plugin插件
	scheduler.AddPlugin(&demo.MockPlugin{})

	// 查看每个task被调度到哪个节点
	time.Sleep(time.Second * 10)
	for _, r := range res {
		r.Exec()
	}

}
