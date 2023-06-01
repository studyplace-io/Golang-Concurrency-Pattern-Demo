package podworker_mode

import (
	"context"
	"k8s.io/klog/v2"
	"time"
)

type Kubelet struct {
	node       string
	podWorkers *PodWorkers
}

func (k *Kubelet) HandlePodAdditions(pods []*Pod) {
	for _, pod := range pods {
		klog.Info(pod.Name)
		k.dispatchWork(pod, SyncPodCreate)
	}
}

func (k *Kubelet) HandlePodUpdates(pods []*Pod) {
	for _, pod := range pods {
		klog.Info(pod.Name)
		k.dispatchWork(pod, SyncPodUpdate)
	}
}

func (k *Kubelet) HandlePodRemoves(pods []*Pod) {
	for _, pod := range pods {
		klog.Info(pod.Name)
		k.dispatchWork(pod, SyncPodKill)
	}
}

func (k *Kubelet) HandlePodSyncs(pods []*Pod) {
	for _, pod := range pods {
		klog.Info(pod.Name)
		k.dispatchWork(pod, SyncPodSync)
	}
}

func (k *Kubelet) dispatchWork(pod *Pod, syncType SyncPodType) {
	// Run the sync in an async worker.
	k.podWorkers.UpdatePod(UpdatePodOptions{
		Pod:        pod,
		UpdateType: syncType,
	})
}

// SyncHandler 处理Pod不同事件的handler
type SyncHandler interface {
	HandlePodAdditions(pods []*Pod)
	HandlePodUpdates(pods []*Pod)
	HandlePodSyncs(pods []*Pod)
	HandlePodRemoves(pods []*Pod)
}

var _ SyncHandler = &Kubelet{}

func (k *Kubelet) Run(updates <-chan PodUpdate) {
	k.syncLoop(context.Background(), updates, k)
}

func (k *Kubelet) syncLoop(ctx context.Context, updates <-chan PodUpdate, handler SyncHandler) {

	syncTicker := time.NewTicker(time.Second * 5)
	defer syncTicker.Stop()
	// 定时清理异常pod的定时器
	// 每两秒检测一次是否有需要清理的 pod
	housekeepingTicker := time.NewTicker(time.Second * 10)
	defer housekeepingTicker.Stop()

	for {
		if !k.syncLoopIteration(ctx, updates, handler, syncTicker.C, housekeepingTicker.C) {
			break
		}
	}
}

func (k *Kubelet) syncLoopIteration(_ context.Context, configCh <-chan PodUpdate, handler SyncHandler,
	syncCh <-chan time.Time, housekeepingCh <-chan time.Time) bool {

	select {
	case u, open := <-configCh:
		if !open {
			klog.ErrorS(nil, "Update channel is closed, exiting the sync loop")
			return false
		}

		switch u.Op {
		case ADD:
			klog.Info("SyncLoop ADD")
			handler.HandlePodAdditions(u.Pods)
		case UPDATE:
			klog.Info("SyncLoop UPDATE")
			handler.HandlePodUpdates(u.Pods)
		case DELETE:
			klog.Info("SyncLoop DELETE")
			handler.HandlePodRemoves(u.Pods)
		default:
			klog.Error(nil, "Invalid operation type received", "operation", u.Op)
		}
	case <-syncCh:
		// 定期上报
		klog.Info("syncCh...")
	case <-housekeepingCh:
		// TODO: 需要定期检测处理坏的容器
		klog.Info("housekeepingCh...")
	}
	return true

}
