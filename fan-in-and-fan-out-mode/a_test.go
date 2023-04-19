package fan_in_and_fan_out_mode

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type User struct {
	Name  string
	Score int64
}

// DoSomethingData 执行
func DoSomethingData(user *User) {
	time.Sleep(time.Second)
	fmt.Println("do something: ", user.Name)
}

func CreatData(num int) []*User {

	out := make([]*User, 0)
	for i := 0; i < num; i++ {
		i := i
		out = append(out, &User{Name: strconv.Itoa(i), Score: int64(i)})
	}
	return out
}

// Handler 类似协程池的思路，固定启动num个goroutine来处理业务(同时读取chan的数据)
func Handler(num int, wg *sync.WaitGroup, s []*User, workerFun func(*User)) {
	inch := make(chan *User, 0)
	// 协程1：把需要处理的参数写入inch
	go func() {
		for _, item := range s {
			inch <- item
		}
		close(inch) // 全部传入就关闭
	}()
	// 协程2：开启num个协程，同时从inch chan中拿数据
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for item := range inch {
				workerFun(item)
			}
		}()
	}
}

// TestOne 不限制goroutine数量，容易会有数量问题
func TestOne(test *testing.T) {
	s := CreatData(100)
	var wg sync.WaitGroup
	// 每个user开启一个协程处理，有几个user对象就会启多少个goroutine
	for _, item := range s {
		wg.Add(1)
		go func(i *User) {
			defer wg.Done()
			DoSomethingData(i)
		}(item)
	}
	wg.Wait()
	//得到数据后下一步处理
	for _, val := range s {
		fmt.Println(val.Score)
	}
}

// TestTwo 限制最大goroutine数量
func TestTwo(test *testing.T) {
	s := CreatData(100)
	var wg sync.WaitGroup
	wg.Add(10)
	Handler(10, &wg, s, DoSomethingData) //只开启特定goroutine去处理
	wg.Wait()
	//得到数据后下一步处理
	for _, val := range s {
		fmt.Println(val.Score)
	}
}

type Receipt struct {
	Name  string
	Score int64
}

// DoData 模拟远程调用数据
func DoData(user *User) *Receipt {
	var res Receipt
	res.Name = user.Name
	res.Score = int64(len(user.Name))
	time.Sleep(time.Second)
	fmt.Println("do something: ", user.Name)
	return &res
}

func HandlerWithReceipt(number int, wg *sync.WaitGroup, s []*User, resCh chan<- *Receipt, workerFun func(*User) *Receipt) {
	inch := make(chan *User, 0)
	//协程1：把需要处理的参数写入inch
	go func() {
		for _, item := range s {
			inch <- item
		}
		close(inch)
	}()
	//协程2：开启number个协程，同时读取参数并把结果写入resCh
	for i := 0; i < number; i++ {
		go func() {
			defer wg.Done()
			for item := range inch {
				res := workerFun(item)
				resCh <- res // 类似FAN-IN，多个goroutine都放入同一个
			}
		}()
	}
}

// TestThree 结合
func TestThree(test *testing.T) {
	s := CreatData(100)
	var wg sync.WaitGroup
	var resWg sync.WaitGroup
	resCh := make(chan *Receipt)
	//协程3：开启一个协程读取结果
	resWg.Add(1)
	go func() {
		defer resWg.Done()
		for item := range resCh {
			fmt.Println(item.Score)
		}
	}()
	wg.Add(10)
	HandlerWithReceipt(10, &wg, s, resCh, DoData)
	wg.Wait()
	close(resCh) //worker结束后需要及时关闭resCh
	resWg.Wait() //保证读取结果完整

}
