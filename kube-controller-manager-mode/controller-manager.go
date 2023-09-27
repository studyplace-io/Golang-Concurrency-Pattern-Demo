package main

import (
	"github.com/practice/Golang-Concurrency-Pattern-Demo/kube-controller-manager-mode/app"
	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register" // for JSON log format registration
	_ "k8s.io/component-base/metrics/prometheus/restclient"
	_ "k8s.io/component-base/metrics/prometheus/version" // for version metric registration
	"os"
)

func main() {
	// 启动kubelet入口
	command := app.NewControllerManagerCommand()
	code := cli.Run(command)
	os.Exit(code)
}
