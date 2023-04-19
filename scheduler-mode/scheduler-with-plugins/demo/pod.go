package demo

import (
	"fmt"
	_interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"
)

type MockPod struct {
	Name          string
	SelectNode    string
	PodRecordNode []_interface.PodRecordNode
}

func (m *MockPod) Exec() {
	fmt.Println("pod", m.Name, "在", m.SelectNode, "节点上")
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
