### 超时退出模式介绍：发送一个请求或是执行一个操作时，同步阻塞等待往往很耗时，因此需要超时退出模式。

- 一句话概括：任何启动并发goroutine的操作，都需要考虑此问题。
- 实现方法：
    1. 使用chan+time.After()来控制超时退出
    2. 使用context+time.After()来控制超时退出
    3. waitGroup上加入Wait()超时退出功能
    4. 超时重试模式
  
- 适用场景：执行一个操作可能会耗时很久的时候使用。
  
```go
1. chan + time.After() 实现超时控制
func TestWaitGroupWithChan(test *testing.T) {
    var wg sync.WaitGroup
    doneC := make(chan struct{})
    
    go func() {
        wg.Wait()
        doneC <- struct{}{}
    }()
    
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            fmt.Println(index)
            // 模拟某个goroutine 执行非常久的时间
            if index == 2 {
                time.Sleep(time.Second * 100)
            }
        }(i)
    }
    
    timeout := time.Duration(10) * time.Second
    fmt.Printf("Wait for waitgroup (up to %s)\n", timeout)
    
    select {
    case <-doneC:
        fmt.Printf("Wait group finished\n")
    case <-time.After(timeout):
        fmt.Printf("Timed out waiting for wait group\n")
    }

}

2. context + time.After() 实现超时控制
func TestContextTimeout2(t *testing.T) {

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()
    stopC := make(chan struct{})
    
    for i := 0; i < 5; i++ {
        go doSomething(ctx, "child goroutine "+strconv.Itoa(i), stopC)
    }
    
    select {
    case <-ctx.Done():
      fmt.Println("call successfully!!!")
      return
    case <-time.After(time.Duration(time.Second * 20)):
      fmt.Println("timeout!!!")
      return
    }

}

```