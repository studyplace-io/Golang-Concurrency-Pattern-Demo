package demo

import (
	_interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/nodes"
	"math/rand"
)

type MockPlugin struct {}

func (m MockPlugin) Score(task _interface.Pod, infos *nodes.NodeInfo) float64 {
	return rand.Float64() * 100
}

func (m MockPlugin) Filter(task _interface.Pod) bool {
	return rand.Float64() != 0
}

var _ _interface.Plugin = &MockPlugin{}


