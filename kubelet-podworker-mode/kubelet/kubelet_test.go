package kubelet

import (
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/kubelet-podworker-mode/kubelet/config"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/kubelet-podworker-mode/kubelet/container"
	"testing"
	"time"
)

func TestKubelet(t *testing.T) {

	k := &Kubelet{
		config:              config.NewKubeletConfig(),
		node:                "test",
		registerNode:        true,
		registerSchedulable: true,
		podManager:          NewBasicPodManager(),
	}

	k.podWorkers = newPodWorkers(k.syncPod, k.syncTerminatedPod, k.syncTerminatingPod)
	e := make(chan PodUpdate, 10)
	// 运行kubelet
	k.Run(e)

	pods := make([]*Pod, 0)

	containers := make([]container.CRI, 0)

	c := &container.Container{
		Name:            "test1-container",
		Image:           "test1-image",
		ImagePullPolicy: "IfNotPresent",
		Status:          container.NoRunning,
	}
	containers = append(containers, c)

	pods = append(pods, &Pod{
		Name:       "test1",
		Containers: containers,
		Status:     NoRunning,
	})
	a := PodUpdate{
		Op:   ADD,
		Pods: pods,
	}
	e <- a

	podsa := []*Pod{}
	containersa := []container.CRI{}

	ca := &container.Container{
		Name:            "test2-container",
		Image:           "test2-image",
		ImagePullPolicy: "IfNotPresent",
		Status:          container.NoRunning,
	}
	containersa = append(containersa, ca)

	podsa = append(podsa, &Pod{
		Name:       "test2",
		Containers: containersa,
		Status:     NoRunning,
	})
	aa := PodUpdate{
		Op:   ADD,
		Pods: podsa,
	}
	e <- aa

	//kube-controller-manager-mode := PodUpdate{
	// Op:   ADD,
	// Pods: podsa,
	//}
	//
	aaa := PodUpdate{
		Op:   DELETE,
		Pods: pods,
	}

	aaaa := PodUpdate{
		Op:   UPDATE,
		Pods: podsa,
	}

	e <- aaa
	e <- aaaa
	//select {}
	time.Sleep(time.Second * 10)

}
