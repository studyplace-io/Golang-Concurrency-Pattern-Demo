package main

import (
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/kubelet-podworker-mode/app"
	"k8s.io/component-base/cli"
	"os"
)

func main() {
	command := app.NewKubeletCommand()
	code := cli.Run(command)
	os.Exit(code)
}
