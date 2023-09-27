package kubelet

import (
	"github.com/practice/Golang-Concurrency-Pattern-Demo/kubelet-podworker-mode/kubelet/container"
)

// Pod 模拟Pod对象
type Pod struct {
	Name       string
	Containers []container.CRI
	Status     PodStatus
}

type PodStatus string

const (
	NoRunning PodStatus = "noRunning"
	Running   PodStatus = "running"
	Failed    PodStatus = "failed"
	UnKnow    PodStatus = "unKnow"
)

// PodOperation 标示Pod的事件
type PodOperation int

const (
	// ADD 新增事件
	ADD PodOperation = iota
	// DELETE 删除事件
	DELETE
	// UPDATE 更新事件
	UPDATE
)

// SyncPodType classifies pod updates, eg: create, update.
type SyncPodType int

const (
	// SyncPodSync is when the pod is synced to ensure desired state
	SyncPodSync SyncPodType = iota
	// SyncPodUpdate is when the pod is updated from source
	SyncPodUpdate
	// SyncPodCreate is when the pod is created from source
	SyncPodCreate
	// SyncPodKill is when the pod should have no running containers. A pod stopped in this way could be
	// restarted in the future due config changes.
	SyncPodKill
)

// PodUpdate
type PodUpdate struct {
	Pods []*Pod
	Op   PodOperation
}

type UpdatePodOptions struct {
	UpdateType SyncPodType
	Pod        *Pod
}
