package scheduler_with_plugins

import (
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/demo"
	_interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/nodes"
	"sync"
	"testing"
	"time"
)

func TestSchedulerWithPlugin(test *testing.T) {

	// 调度器
	scheduler := Scheduler{
		pods: make(chan _interface.Pod, 10),
		nodeInfos: nodes.NewNodeInfos(),
		workers: 8,
		plugins: make([]_interface.Plugin, 0),
		wg: sync.WaitGroup{},
	}

	// 加入模拟node
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node1"))
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node2"))
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node3"))

	// 启动调度器
	scheduler.wg.Add(1)
	go scheduler.run()

	go func() {
		scheduler.wg.Wait()
	}()

	res := make([]*demo.MockPod, 0)

	t1 := &demo.MockPod{
		Name: "pod1",
		PodRecordNode: make([]_interface.PodRecordNode, 0),
	}
	scheduler.AddPod(t1)
	res = append(res, t1)

	t2 := &demo.MockPod{
		Name: "pod1",
		PodRecordNode: make([]_interface.PodRecordNode, 0),
	}
	scheduler.AddPod(t2)
	res = append(res, t2)

	t3 := &demo.MockPod{
		Name: "pod1",
		PodRecordNode: make([]_interface.PodRecordNode, 0),
	}
	scheduler.AddPod(t3)
	res = append(res, t3)

	scheduler.AddPlugin(&demo.MockPlugin{})

	// 查看每个task被调度到哪个节点
	time.Sleep(time.Second * 10)
	for _, r := range res {
		r.Exec()
	}

}
