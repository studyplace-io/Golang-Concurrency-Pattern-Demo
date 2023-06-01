package podworker_mode

import (
	"context"
	container2 "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/kubelet-podworker-mode/container"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	"sync"
)

type PodWorkers struct {
	podUpdates           map[string]chan PodWork
	podLock              sync.Mutex
	syncPodFn            syncPodFnType
	syncTerminatingPodFn syncTerminatingPodFnType
	syncTerminatedPodFn  syncTerminatedPodFnType
}

func newPodWorkers(syncPodFn syncPodFnType, syncTerminatingPodFn syncTerminatingPodFnType, syncTerminatedPodFn syncTerminatedPodFnType) *PodWorkers {
	return &PodWorkers{
		podUpdates:           make(map[string]chan PodWork),
		syncPodFn:            syncPodFn,
		syncTerminatingPodFn: syncTerminatingPodFn,
		syncTerminatedPodFn:  syncTerminatedPodFn,
	}
}

func (k *Kubelet) syncPod(_ context.Context, pod *Pod) error {
	klog.Info("syncPod....")
	// TODO: 这里可以对容器进行操作
	// 1. 首先如果状态不是running，代表第一次启动
	// 2. 遍历启动每个容器，并且改容器状态为running，如果有任何个失败，退出，并改成failed
	// 3. 如果都正常running，就把pod改为running

	// 如果pod原本就是running，遍历检查容器状态，如果出现fail，就修改pod的状态为failed

	if pod.Status == NoRunning {
		for _, container := range pod.Containers {

			c := container.(*container2.Container)

			err := c.PullImage(c.Image)
			if err != nil {
				c.Status = container2.Fail
				pod.Status = Failed
				return nil
			}
			err = c.RunPodSandbox()
			if err != nil {
				c.Status = container2.Fail
				pod.Status = Failed
				return nil
			}
			err = c.CreateContainer()
			if err != nil {
				// TODO: 这里要退出容器操作
				c.Status = container2.Fail
				pod.Status = Failed
				return nil
			}
			err = c.StartContainer()
			if err != nil {
				c.Status = container2.Fail
				pod.Status = Failed
				return nil
			}
			c.Status = container2.Running
		}
		pod.Status = Running
	}

	if pod.Status == Running {
		for _, container := range pod.Containers {
			c := container.(*container2.Container)
			if c.Status != container2.Running {
				pod.Status = Failed
				// TODO: 这里需要考虑容器退出处理
			}
		}
	}

	return nil
}

func (k *Kubelet) syncTerminatingPod(_ context.Context, pod *Pod) error {
	klog.Info("syncTerminatingPod....")
	// 如果pod是running状态，调用CRI接口停止容器
	// 如果是fail状态，直接退出

	return nil
}

func (k *Kubelet) syncTerminatedPod(_ context.Context, pod *Pod) error {
	klog.Info("syncTerminatedPod....")
	return nil
}

type syncPodFnType func(ctx context.Context, pod *Pod) error
type syncTerminatingPodFnType func(ctx context.Context, pod *Pod) error
type syncTerminatedPodFnType func(ctx context.Context, pod *Pod) error

// PodWorkType classifies the three phases of pod lifecycle - setup (sync),
// teardown of containers (terminating), cleanup (terminated).
type PodWorkType int

const (
	// SyncPodWork is when the pod is expected to be started and running.
	SyncPodWork PodWorkType = iota
	// TerminatingPodWork is when the pod is no longer being set up, but some
	// containers may be running and are being torn down.
	TerminatingPodWork
	// TerminatedPodWork indicates the pod is stopped, can have no more running
	// containers, and any foreground cleanup can be executed.
	TerminatedPodWork
)

type PodWork struct {
	WorkType PodWorkType
	Pod      *Pod
}

func (p *PodWorkers) UpdatePod(options UpdatePodOptions) {

	pod := options.Pod
	p.podLock.Lock()
	defer p.podLock.Unlock()

	podUpdates, exists := p.podUpdates[options.Pod.Name]
	if !exists {
		podUpdates = make(chan PodWork, 1)
		p.podUpdates[options.Pod.Name] = podUpdates

		go func() {
			klog.Infof("pod %s start a goroutine to handle events...", pod.Name)
			defer func() {
				klog.Info("one pod worker exit....")
			}()
			defer runtime.HandleCrash()
			p.managePodLoop(podUpdates)
		}()
		podUpdates <- PodWork{
			WorkType: SyncPodWork,
			Pod:      pod,
		}

	} else {
		klog.Infof("pod %s has already started a goroutine to handle events...", pod.Name)
		if options.UpdateType == SyncPodUpdate || options.UpdateType == SyncPodSync {
			podUpdates <- PodWork{
				WorkType: SyncPodWork,
				Pod:      pod,
			}
		}

		if options.UpdateType == SyncPodKill {
			podUpdates <- PodWork{
				WorkType: TerminatingPodWork,
				Pod:      pod,
			}
		}
	}

}

func (p *PodWorkers) managePodLoop(podUpdates <-chan PodWork) {

	for update := range podUpdates {

		err := func() error {
			var err error
			switch {
			case update.WorkType == TerminatedPodWork:
				err = p.syncTerminatedPodFn(context.Background(), update.Pod)
			case update.WorkType == TerminatingPodWork:
				err = p.syncTerminatingPodFn(context.Background(), update.Pod)
			default:
				err = p.syncPodFn(context.Background(), update.Pod)
			}
			return err
		}()

		switch {
		case err == context.Canceled:
			klog.Info("the context is canceled...")
		case err != nil:
			klog.Error("pod sync err: ", err)
		case update.WorkType == TerminatedPodWork:
			klog.Info("Processing pod event done", " pod ", update.Pod.Name, " updateType ", update.WorkType)
			return
		case update.WorkType == TerminatingPodWork:
			klog.Info("Processing pod event done", " pod ", update.Pod.Name, " updateType ", update.WorkType)
			return
		default:
			klog.Info("pod sync success...")
		}

	}

}
