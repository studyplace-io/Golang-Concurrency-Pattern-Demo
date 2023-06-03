package kubelet

import (
	"sync"
)

// Manager 管理器
type Manager interface {
	// GetPods returns the regular pods bound to the kubelet and their spec.
	GetPods() []*Pod
	GetPodByName(name string) (*Pod, bool)
	// SetPods replaces the internal pods with the new pods.
	// It is currently only used for testing.
	SetPods()
	// AddPod adds the given pod to the manager.
	AddPod(pod *Pod)
	// UpdatePod updates the given pod in the manager.
	UpdatePod(pod *Pod)
	// DeletePod deletes the given pod from the manager.  For mirror pods,
	// this means deleting the mappings related to mirror pods.  For non-
	// mirror pods, this means deleting from indexes for all non-mirror pods.
	DeletePod(pod *Pod)
}

type basicManager struct {
	// Protects all internal maps.
	lock sync.RWMutex
	// Pods indexed by full name for easy access.
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

// Set the internal pods based on the new pods.
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
