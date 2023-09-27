package _interface

import "github.com/practice/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/nodes"

// Pod 接口
type Pod interface {
	Exec()
	SetNode(string)
	SetPodRecordNode(string, float64)
	GetPodRecordNodeList() []PodRecordNode
}

// PodRecordNode 为了记录特定pod在所有node中的打分
type PodRecordNode struct {
	NodeName string
	Score    float64
}

// Plugin 插件接口
type Plugin interface {
	Score(pod Pod, infos *nodes.NodeInfo) float64
	Filter(pod Pod) bool
}
