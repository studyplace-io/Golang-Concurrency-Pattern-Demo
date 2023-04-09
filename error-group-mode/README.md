### ErrorGroup模式介绍：想要并发运行业务时会直接开goroutine，但是直接go函数是无法对返回数据进行处理error，可以使用ErrorGroup。

- 一句话概括：goroutine出去的方法，可以返回需要的error。
- 实现方法：有两种实现方式：
    1. waitgroup + chan error收集error。
    2. 使用ErrorGroup包
- 适用场景：并发goroutine需要返回error时适用。
  1. goroutine需要返回error时使用。
  2. 整匹goroutine只要一个报错就不再执行时使用。
  
```go
1. waitgroup+chan error

type ErrorGroupExecutor struct {
    errC chan error
    wg   sync.WaitGroup
}

// Execute 执行：如果报错，把错误放入errC中
func (e *ErrorGroupExecutor) Execute(num int)  {
    defer e.wg.Done()
    
    if num%2 == 0{
        str := "err" + strconv.Itoa(num)
        e.errC <-errors.New(str)
    }
}

2. 使用ErrorGroup包

func TryUseErrGroup() {

    var urls = []string{
        "http://www.golang.org/",
        "http://www.baidu.com/",
        "http://www.noexist11111111.com/",
        "http://www.noexist11111111.com/",
        "http://www.noexist11111111.com/",
        "http://www.noexist11111111.com/",
    }

    g := new(errgroup.Group)

    for _, url := range urls {
        url := url
        g.Go(func() error {
            resp, err := http.Get(url)
            if err != nil {
                return err
            }
            fmt.Printf("get [%s] success: [%d] \n", url, resp.StatusCode)
            return resp.Body.Close()

         })
    }

    // 注意：这里的wait是 "当所有的goroutine中，只要有一个报错，就会返回且退出"
    if err := g.Wait(); err != nil {
        // 只会返回"第一个error"，不管后续还有没有error，都不再执行
        fmt.Println(err)
    } else {
        fmt.Println("All success!")
    }
    fmt.Println("主goroutine阻塞结束，进程退出")
}
```