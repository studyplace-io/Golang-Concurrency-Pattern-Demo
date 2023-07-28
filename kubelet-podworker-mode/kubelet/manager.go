package kubelet

import (
	"sync"
)

// Manager 管理器
type Manager interface {
	// GetPods 获取所有的pod
	GetPods() []*Pod
	// GetPodByName 获取特定pod
	GetPodByName(name string) (*Pod, bool)
	// SetPods 初始化
	SetPods()
	// AddPod 把传入的pod加入的管理器中
	AddPod(pod *Pod)
	// UpdatePod 更新管理器中的pod
	UpdatePod(pod *Pod)
	// DeletePod 删除管理器中的pod
	DeletePod(pod *Pod)
}

// basicManager pod的本地缓存
type basicManager struct {
	lock sync.RWMutex
	// podByFullName 用来存储pod对象
	podByFullName map[string]*Pod
}

// NewBasicPodManager returns a functional Manager.
func NewBasicPodManager() Manager {
	pm := &basicManager{
		lock: sync.RWMutex{},
	}
	pm.SetPods()
	return pm
}

func (pm *basicManager) SetPods() {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.podByFullName = make(map[string]*Pod)

}

func (pm *basicManager) AddPod(pod *Pod) {
	pm.UpdatePod(pod)
}

func (pm *basicManager) UpdatePod(pod *Pod) {
	pm.updatePodsInternal(pod)
}

func (pm *basicManager) updatePodsInternal(pods ...*Pod) {
	for _, pod := range pods {
		podFullName, ok := pm.GetPodByName(pod.Name)
		// This logic relies on a static pod and its mirror to have the same name.
		// It is safe to type convert here due to the IsMirrorPod guard.
		pm.lock.Lock()
		if ok {
			pm.podByFullName[podFullName.Name] = pod
		} else {
			pm.podByFullName[pod.Name] = pod
		}
		pm.lock.Unlock()
	}
}

func (pm *basicManager) DeletePod(pod *Pod) {
	podFullName, ok := pm.GetPodByName(pod.Name)
	if ok {
		pm.lock.Lock()
		defer pm.lock.Unlock()
		delete(pm.podByFullName, podFullName.Name)
	}
}

func (pm *basicManager) GetPods() []*Pod {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	return podsMapToPods(pm.podByFullName)
}

func podsMapToPods(UIDMap map[string]*Pod) []*Pod {
	pods := make([]*Pod, 0, len(UIDMap))
	for _, pod := range UIDMap {
		pods = append(pods, pod)
	}
	return pods
}

func (pm *basicManager) GetPodByName(name string) (*Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByFullName[name]
	return pod, ok
}
