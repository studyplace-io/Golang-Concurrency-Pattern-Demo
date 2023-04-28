package demo

import (
	_interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"
	"k8s.io/klog/v2"
)

type MockPod struct {
	Name          string
	SelectNode    string
	PodRecordNode []_interface.PodRecordNode
}

func (m *MockPod) Exec() {
	klog.Info("pod ", m.Name, " in ", m.SelectNode, " node")
}

func (m *MockPod) SetNode(s string) {
	m.SelectNode = s
}

func (m *MockPod) SetPodRecordNode(s string, f float64) {
	r := _interface.PodRecordNode{
		Score:    f,
		NodeName: s,
	}
	m.PodRecordNode = append(m.PodRecordNode, r)
}

func (m *MockPod) GetPodRecordNodeList() []_interface.PodRecordNode {
	return m.PodRecordNode
}

var _ _interface.Pod = &MockPod{}
