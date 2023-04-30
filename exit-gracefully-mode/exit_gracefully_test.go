package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"time"
)

func Example() {

	// 优雅退出
	// 方式一：
	exitSignalHandler()

	// 方式二：
	SetupSignalHandler(shutdown)

}

// exitSignalHandler 接收退出信号
func exitSignalHandler() {
	// 接受一个signal信号的channel
	signalChan := make(chan os.Signal, 1)

	// 使用os库中的Notify 接受os传来的信号，放到chan中
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// 采用一个goroutine执行业务逻辑，多是一个常驻goroutine
	go func() {
		fmt.Println("do something.....")
	}()

	// 阻塞
	for {
		select {
		// 如果没有收到退出信号，会一直阻塞再这里，或是使用default不让其阻塞
		// 收到退出信号会return
		case sig := <-signalChan:
			log.Printf("收到进程给的信号[signal = %v ]", sig)
			return
		default:
			log.Println("尚未收到信号")
			time.Sleep(time.Second)
		}
	}
}

// SetupSignalHandler 阻塞接收退出信号，并输入一个shutdownFunc，执行退出前的资源清理动作。
func SetupSignalHandler(shutdownFunc func(bool)) {
	// 接受os通知的chan
	closeSignalChan := make(chan os.Signal, 1)
	// 信号通知函数
	signal.Notify(closeSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	// 启一个goroutine执行业务逻辑
	go func() {
		fmt.Println("do something.....")
	}()

	// 当没有信号传入时，会阻塞
	sig := <-closeSignalChan
	log.Printf("收到信号[signal = %v ]", sig)
	// 调用shutdownFunc前，判断退出信号的种类
	shutdownFunc(sig == syscall.SIGQUIT)

}

// shutdown 能够区分强制退出与优雅退出的退出逻辑，执行不同的动作
func shutdown(isGraceful bool) {

	// 查看通知种类，判断是否优雅退出

	if isGraceful {
		fmt.Println("graceful exit...")
		//当满足 sig == syscall.SIGQUIT,做相应退出处理 ex: 清理资源
	}
	fmt.Println("force exit...")
	// 不是syscall.SIGQUIT的退出信号时，做相应退出处理

}
