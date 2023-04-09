### 优雅退出模式介绍：生产环境下运行的系统几乎都需要优雅退出,即程序接收退出通知后,会有机会先执行一段清理代码,将收尾工作做完后再真正退出。

- 一句话概括：在代码关停或是crash时，能够做一些逻辑操作。
- 实现方法：
    1. 使用os.Signal的chan，并使用Notify()进行通知。
- 适用场景：服务启动时，都需要考虑"优雅退出"。
  1. server启动需要优雅退出。
  
```go

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
		fmt.Println("这里执行业务逻辑。。。。。")
	}()
	
        // 这里会阻塞
	sig := <-closeSignalChan
	log.Printf("收到信号[signal = %v ]", sig)
	// 调用shutdownFunc前，判断退出信号的种类
	shutdownFunc(sig == syscall.SIGQUIT)

}

```