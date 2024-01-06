package tree_template

import (
	"fmt"
	"github.com/study-io/Golang-Concurrency-Pattern-Demo/tree-template-engine/model"
	"log"
	"os"
	"os/exec"

	"sync"
	"syscall"
)

// runCommand 执行 bash 命令
func runCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		// 如果命令执行失败，则返回错误
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			// 获取命令的退出状态码
			status := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
			return fmt.Errorf("command failed with exit code %d", status)
		}
		return err
	}
	return nil
}

// handleFailure 处理错误
func handleFailure(task *model.Task, status *model.Status) {
	switch task.OnFailure {
	case "stop":
		os.Exit(1)
	case "skip":
		return
	case "retry":
		task.Failed = false
		task.Completed = false
		runTask(task, status)
	default:
		err := runCommand(task.OnFailure)
		if err != nil {
			log.Printf("Error executing custom script: %v\n", err)
		}
	}
}

// runTask 递归执行命令
func runTask(task *model.Task, status *model.Status) {

	switch task.Type {
	// 执行命令
	case "command":
		for _, command := range task.Commands {
			err := runCommand(command)
			if err != nil {
				log.Printf("Error executing command: %v\n", err)
				task.Failed = true          // 标记任务为失败状态
				task.Completed = true       // 完成任务以跳过 on_failure 任务
				handleFailure(task, status) // 执行 on_failure 任务
				return
			}
		}
	// 串行执行
	case "serial":
		for _, child := range task.Children {
			runTask(&child, status)
			if !child.Completed {
				task.Failed = true
			}
		}
	// 并行执行
	case "parallel":
		var wg sync.WaitGroup
		wg.Add(len(task.Children))
		for i := range task.Children {
			go func(child model.Task) {
				runTask(&child, status)
				wg.Done()
			}(task.Children[i])
		}
		wg.Wait()

		task.Failed = false
		for _, child := range task.Children {
			if child.Failed {
				task.Failed = true
				break
			}
		}
	}

	task.Completed = true
	status.Lock()
	defer status.Unlock()
	for _, t := range status.Tasks {
		if !t.Completed {
			return
		}
	}
}
