### 定时任务 Worker-Job 模式介绍：

- 一句话概括：定时任务可使用
- 实现方法：
    1. 抽象出 Worker 接口
    2. 用户输入 func 业务逻辑，使用 k8s 工具包实现定时轮循任务