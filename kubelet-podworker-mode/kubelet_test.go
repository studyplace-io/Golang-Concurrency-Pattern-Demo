package podworker_mode

import (
	"testing"
	"time"
)

func TestKubelet(t *testing.T) {

	k := &Kubelet{
		node: "test",
	}
	k.podWorkers = newPodWorkers(k.syncPod, k.syncTerminatedPod, k.syncTerminatingPod)
	e := make(chan PodUpdate, 10)
	go k.Run(e)

	pods := []*Pod{}
	pods = append(pods, &Pod{
		Name: "ddd",
	})

	podsa := []*Pod{}
	podsa = append(podsa, &Pod{
		Name: "ddddcccdd",
	})

	a := PodUpdate{
		Op:   ADD,
		Pods: pods,
	}

	e <- a

	aa := PodUpdate{
		Op:   ADD,
		Pods: podsa,
	}

	aaa := PodUpdate{
		Op:   DELETE,
		Pods: pods,
	}

	aaaa := PodUpdate{
		Op:   UPDATE,
		Pods: podsa,
	}

	e <- aa
	e <- aaa
	e <- aaaa

	time.Sleep(time.Second * 5)

}
