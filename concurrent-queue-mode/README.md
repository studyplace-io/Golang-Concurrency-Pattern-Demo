### 并发队列模式介绍：支持并发操作的队列

- 一句话概括：使用sync.Cond实现并发通知
- 实现方法：
    1. 实现Queue接口，基础队列功能
    2. 实现ConcurrentQueue接口，使用sync.Cond实现并发goroutine的阻塞与通知功能