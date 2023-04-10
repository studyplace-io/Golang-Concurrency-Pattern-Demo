### 流水线模式介绍：由好几个不同阶段组成(方法)，每个方法input与output都使用"chan"来传递。

- 一句话概括：可以连贯下去的流程(第一步干啥，第二步干啥等)，并当中有数据的参与时使用。
- 实现方法：
    1. 参与流水线的方法的input与output都使用chan来传递(类似消费者-生产者)。
  
      1. 每个阶段把数据通过chan传递给下一个阶段。
      2. 每个阶段要创建1个goroutine和1个chan，此goroutine向chan写数据，函数要返回这个chan。
      3. 用1个函数来组织流水线，例如main函数。
  
- 适用场景：当流程有步骤化时可以使用，在并发场景下可以方便进行扩展。
  
```go
// producer 负责生产数据，返回一个chan，把准备好的数据放入chan中。
func producer(num ...int) <-chan int {
    out := make(chan int)
    
    // 异步启goroutine准备数据，并放入chan
    go func() {
        defer close(out)
        for _, n := range num {
            out <- n
        }
    }()
    
    return out
}

// square 执行主要的业务逻辑，从准备好的chan中拿取数据，并执行业务逻辑，执行后放入chan中
func square(inputC <-chan int) <-chan int {

    out := make(chan int)
    // 异步启goroutine准备数据，并放入chan
    go func() {
        defer close(out)
        for n := range inputC {
            out <- n * n
        }
    }()
    return out
}

func TestPipeline(t *testing.T) {
    // 流水线模式
    in := producer(1, 2, 3, 4)
    ch := square(in)
    
    // 消费准备好的数据，ex: 落库等操作
    for ret := range ch {
        fmt.Println("res: ", ret)
    }
    
    fmt.Println("pipeline mode finished.")
}

```