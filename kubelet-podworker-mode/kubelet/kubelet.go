package kubelet

import (
	"context"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/kubelet-podworker-mode/kubelet/config"
	"k8s.io/klog/v2"
	"time"
)

// Kubelet 对象
type Kubelet struct {
	// config 配置文件
	config *config.KubeletConfig
	// node 节点名，记录本kubelet运行的节点
	node string
	// podWorkers 工作队列，负责对目前kubelet的pod进行生命周期管理
	podWorkers PodWorkers
	// podManager pod记录管理器，负责记录目前kubelet运行的所有pod
	podManager Manager

	// Set to true to have the node register itself with the apiserver.
	registerNode bool
	// Set to true to have the node register itself as schedulable.
	registerSchedulable bool
}

func (k *Kubelet) HandlePodCleanups(ctx context.Context) error {
	klog.Info("HandlePodCleanups...")
	allPods := k.podManager.GetPods()

	failPods := make([]*Pod, 0)

	for _, pod := range allPods {
		if pod.Status == Failed {
			failPods = append(failPods, pod)
		}
	}

	for _, pod := range failPods {
		err := k.syncTerminatingPod(ctx, pod)
		k.podManager.DeletePod(pod)
		if err != nil {
			return err
		}
	}

	return nil

}

func (k *Kubelet) HandlePodAdditions(pods []*Pod) {
	for _, pod := range pods {
		klog.Info(pod.Name)
		k.dispatchWork(pod, SyncPodCreate)
		k.podManager.AddPod(pod)
	}
}

func (k *Kubelet) HandlePodUpdates(pods []*Pod) {
	for _, pod := range pods {
		klog.Info(pod.Name)
		k.dispatchWork(pod, SyncPodUpdate)
		k.podManager.UpdatePod(pod)
	}
}

func (k *Kubelet) HandlePodRemoves(pods []*Pod) {
	for _, pod := range pods {
		klog.Info(pod.Name)
		k.dispatchWork(pod, SyncPodKill)
		k.podManager.DeletePod(pod)
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
	// HandlePodAdditions 处理pod新增事件
	HandlePodAdditions(pods []*Pod)
	// HandlePodUpdates 处理pod更新事件
	HandlePodUpdates(pods []*Pod)
	// HandlePodSyncs 处理pod同步事件
	HandlePodSyncs(pods []*Pod)
	// HandlePodRemoves 处理pod删除事件
	HandlePodRemoves(pods []*Pod)
	// HandlePodCleanups 处理pod坏死事件
	HandlePodCleanups(ctx context.Context) error
}

var _ SyncHandler = &Kubelet{}

// Run 执行
func (k *Kubelet) Run(updates <-chan PodUpdate) {
	// 循环事件
	go k.syncLoop(context.Background(), updates, k)
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

func (k *Kubelet) syncLoopIteration(ctx context.Context, configCh <-chan PodUpdate, handler SyncHandler,
	syncCh <-chan time.Time, housekeepingCh <-chan time.Time) bool {

	select {
	case u, open := <-configCh:
		if !open {
			klog.ErrorS(nil, "Update channel is closed, exiting the sync loop")
			return false
		}

		// 区分不同事件
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
		// TODO: 定期上报容器状态
		klog.Info("syncCh...")
	case <-housekeepingCh:
		// TODO: 需要定期检测处理坏的容器
		klog.Info("housekeepingCh...")
		if err := handler.HandlePodCleanups(ctx); err != nil {
			klog.ErrorS(err, "Failed cleaning pods")
		}
	}
	return true

}
