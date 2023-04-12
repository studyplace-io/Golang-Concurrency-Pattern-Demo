### Future模式介绍：当“未来”需要的结果（一般就是一个网络请求的结果），现在就要发起请求或操作时使用的模式。

- 一句话概括：在后台执行一个异步请求，等以后的某个时间要调用时可以拿到。
- 实现方法：
    1. 在func中启动goroutine且返回chan接口的方式来实现
  
- 适用场景：执行一个操作可能会耗时很久的时候使用，主流程中可以执行其他操作。
  
```go
type query struct {
    // 查询参数 chan
    sql chan string
    // 接收结果参数
    result chan string
}

// execQuery 执行查询db任务
func execQuery(q *query) {
    go func() {
        queryCmd := <-q.sql
        fmt.Println("查询db，耗时任务")
        time.Sleep(time.Second * 10)
        q.result <- "result from " + queryCmd
    }()
}

func TestFutureMode(test *testing.T) {
    q := newQuery()
    
    go execQuery(q)
    q.sql <- "select * from table"
    time.Sleep(10 * time.Second)
    
    fmt.Println("我这里还能做好多事情。。。。。")
    
    fmt.Println(<-q.result)
    
    q.sql <- "select * from table aaa "
    time.Sleep(10 * time.Second)
    fmt.Println(<-q.result)
}


```