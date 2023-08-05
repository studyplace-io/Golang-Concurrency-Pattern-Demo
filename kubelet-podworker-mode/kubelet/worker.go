package kubelet

import (
	"context"
	"errors"
	container2 "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/kubelet-podworker-mode/kubelet/container"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	"sync"
)

// PodWorkers kubelet中主要干活的组件
type PodWorkers interface {
	// 所有handler都走此方法
	UpdatePod(options UpdatePodOptions)
}

// podWorkers 负责所有pod的生命周期
type podWorkers struct {
	// 管理每个pod的map
	podUpdates map[string]chan PodWork
	podLock    sync.Mutex
	// pod的同步事件处理func
	syncPodFn            syncPodFnType
	syncTerminatingPodFn syncTerminatingPodFnType
	syncTerminatedPodFn  syncTerminatedPodFnType
}

func newPodWorkers(syncPodFn syncPodFnType, syncTerminatingPodFn syncTerminatingPodFnType, syncTerminatedPodFn syncTerminatedPodFnType) *podWorkers {
	return &podWorkers{
		podUpdates:           make(map[string]chan PodWork),
		syncPodFn:            syncPodFn,
		syncTerminatingPodFn: syncTerminatingPodFn,
		syncTerminatedPodFn:  syncTerminatedPodFn,
	}
}

// preCheckForRunContainer 检查此节点是否可以运行
func (k *Kubelet) preCheckForRunContainer(pod *Pod) error {
	ok := k.canRunPod(pod)
	if !ok {
		return errors.New("check security permissions fail")
	}

	err := k.buildCgroups(pod)
	if err != nil {
		return err
	}
	err = k.makePodDataDirs(pod)
	if err != nil {
		return err
	}
	return err
}

// canRunPod 判断此node是否可以运行此pod，确保 pod 有正确的安全权限
func (k *Kubelet) canRunPod(pod *Pod) bool {
	klog.Info("check security permissions...")
	return true
}

// buildCgroups 为pod的容器设置Cgroups
func (k *Kubelet) buildCgroups(pod *Pod) error {
	klog.Info("Set up cgroups for all containers of the pod...")
	return nil
}

// makePodDataDirs 为pod的容器挂载需要的目录
func (k *Kubelet) makePodDataDirs(pod *Pod) error {
	klog.Info("Create the required directories for the pod's containers...")
	return nil
}

func (k *Kubelet) syncPod(_ context.Context, pod *Pod) error {
	klog.Info("syncPod....")
	// 1. 首先如果状态不是running，代表第一次启动
	// 2. 遍历启动每个容器，并且改容器状态为running，如果有任何个失败，退出，并改成failed
	// 3. 如果都正常running，就把pod改为running
	if pod.Status == NoRunning {
		// 检查是否能运行容器
		err := k.preCheckForRunContainer(pod)
		if err != nil {
			pod.Status = Failed
			return nil
		}

		for _, container := range pod.Containers {

			c := container.(*container2.Container)

			err := c.PullImage(c.Image)
			if err != nil {
				pod.Status = Failed
				return nil
			}
			err = c.RunPodSandbox()
			if err != nil {
				pod.Status = Failed
				return nil
			}
			err = c.CreateContainer()
			if err != nil {
				// TODO: 这里要退出容器操作
				pod.Status = Failed
				return nil
			}
			err = c.StartContainer()
			if err != nil {
				pod.Status = Failed
				return nil
			}
		}
		pod.Status = Running
		k.podManager.AddPod(pod)
	}

	// 如果pod原本就是running，遍历检查容器状态，如果出现fail，就修改pod的状态为failed
	if pod.Status == Running {
		for _, container := range pod.Containers {
			c := container.(*container2.Container)
			if c.Status != container2.Running {
				pod.Status = Failed
				// FIXME: 这里需要考虑容器退出处理
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

// PodWork podWorker管理的对象，也就是pod外面再包装一层
// 记录本次的pod事件类型与pod对象本身
type PodWork struct {
	WorkType PodWorkType
	Pod      *Pod
}

// UpdatePod 管理pod的主要逻辑
func (p *podWorkers) UpdatePod(options UpdatePodOptions) {
	pod := options.Pod
	p.podLock.Lock()
	defer p.podLock.Unlock()

	// 查map，是否已经启动goroutine
	podUpdates, exists := p.podUpdates[options.Pod.Name]
	if !exists {
		podUpdates = make(chan PodWork, 1)
		p.podUpdates[options.Pod.Name] = podUpdates
		// 第一次使用pod worker管理一个pod的生命周期，需要启动goroutine管理
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

func (p *podWorkers) managePodLoop(podUpdates <-chan PodWork) {

	// 不断从chan中取出pod
	for update := range podUpdates {

		// 区分不同命令的pod
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

		// 不同结果
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
