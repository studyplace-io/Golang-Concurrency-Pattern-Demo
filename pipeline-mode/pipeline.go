package pipeline_mode

/*
 流水线的特点:
 1. 每个阶段把数据通过channel传递给下一个阶段。
 2. 每个阶段要创建1个goroutine和1个通道，这个goroutine向里面写数据，函数要返回这个通道。
 3. 有1个函数来组织流水线，我们例子中是main函数。
*/

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
