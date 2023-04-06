package error_group_mode

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
)

/*
	实现类似errgroup的功能
	初始化一个chan error类型的chan，然后每个goroutine在遇到错误时候，
	将error写入chan，这样主goroutine通过for range去遍历这个chan就行了。
*/

var wg sync.WaitGroup

func TestErrorPractice(t *testing.T) {
	num := 5
	TryUseChanAndErrorGroup(num)
}

func TryUseChanAndErrorGroup(num int) {


	errChan := make(chan error, num)
	wg.Add(num)
	for i :=0 ; i < num; i++ {
		// 模拟goroutine出错，使用chan 捞出error结果
		go func(i int) {
			defer wg.Done()
			str := "err" + strconv.Itoa(i)
			errChan <- errors.New(str)
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Println(err)
	}

	fmt.Println("主goroutine退出")

}
