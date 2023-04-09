package error_group_mode

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
)

/*
	实现类似errgroup的功能:
	初始化一个chan error类型的chan，然后每个goroutine在遇到错误时候，
	将error写入chan，这样主goroutine通过for range去遍历这个chan。
*/


func TestErrorPractice(t *testing.T) {
	num := 10
	TryUseErrorChan(num)
}

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

func NewExecutor(errC chan error) *ErrorGroupExecutor {
	return &ErrorGroupExecutor{
		errC: errC,
		wg: sync.WaitGroup{},
	}
}

func TryUseErrorChan(num int) {
	errChan := make(chan error, num)
	e := NewExecutor(errChan)

	for i := 0; i < num; i++ {
		e.wg.Add(1)
		go e.Execute(i)
	}

	e.wg.Wait()		// 阻塞等待
	close(e.errC)	// 多生产者，需要关闭chan

	// 消费查看errC
	for err := range e.errC {
		fmt.Println(err)
	}

	fmt.Println("主goroutine退出")


}
