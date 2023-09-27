package app

import (
	"context"
	"fmt"
	"github.com/practice/Golang-Concurrency-Pattern-Demo/kube-controller-manager-mode/pkg/controller/mycontroller1"
	"github.com/practice/Golang-Concurrency-Pattern-Demo/kube-controller-manager-mode/pkg/controller/mycontroller2"
	"k8s.io/klog/v2"
)

// startMyController1 初始化deployment controller
func startMyController1(ctx context.Context) (Interface, bool, error) {
	klog.Infof("start myController1...")
	mc, err := mycontroller1.NewController()
	if err != nil {
		return nil, true, fmt.Errorf("error creating Deployment controller: %v", err)
	}
	// 异步启动controller
	go mc.Run(ctx, mycontroller1.DefaultWorkers)
	return nil, true, nil
}

// startMyController2 初始化deployment controller
func startMyController2(ctx context.Context) (Interface, bool, error) {
	klog.Infof("start myController2...")
	mc, err := mycontroller2.NewController()
	if err != nil {
		return nil, true, fmt.Errorf("error creating Deployment controller: %v", err)
	}
	// 异步启动controller
	go mc.Run(ctx, mycontroller1.DefaultWorkers)
	return nil, true, nil
}
