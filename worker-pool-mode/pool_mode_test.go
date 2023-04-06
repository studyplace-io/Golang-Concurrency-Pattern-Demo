package worker_pool_mode

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	Running = "running"
	NoUSE   = "nouse"
)

// Obj 需要放入池中的对象
type Obj struct {
	Id	   int
	Status string
	// 这里可以加入需要的字段或方法
}

// 自定义Obj的方法

func (o Obj) GetId() int {
	return o.Id
}

func (o Obj) Execute() {
	o.Status = Running // 转变状态
}


// ObjPool 池
type ObjPool struct {
	pool   	   chan *Obj   // 维护一个保存对象的chan
	poolSize   int		   // pool 大小
	poolInit   bool        // 标示是否已初始化
}

// NewObjPool 创建新的pool, 并初始化中的对象
func NewObjPool(taskNum int, init bool) *ObjPool {
	taskPool := &ObjPool{
		pool: make(chan *Obj, taskNum),
		poolSize: taskNum,
		poolInit: init,
	}

	// 初始化pool对象
	if taskPool.poolInit {
		for i := 0; i < taskPool.poolSize; i++ {
			taskPool.pool <- &Obj{
				Id: i,
				Status: NoUSE,
			}
		}
	}

	return taskPool
}

var once sync.Once

// GetObj 从池中取出对象
func (t *ObjPool) GetObj(timeout time.Duration) (*Obj, error) {
	// 以防万一，再次执行执行init
	poolInit := func() {
		for i := 0; i < t.poolSize; i++ {
			t.pool <- &Obj{
				Id: i,
				Status: NoUSE,
			}
		}
	}

	if ! t.poolInit {
		once.Do(poolInit)
		t.poolInit = true
	}

	// 使用select取出池中obj, 或是超时返回
	select {
	case obj := <-t.pool:
		return obj, nil
	case <-time.After(timeout):
		return nil, errors.New("get obj from pool timeout error")
	}

}

// PutObj 将obj放回池中
func (t *ObjPool) PutObj(obj *Obj) error {
	// FIXME: 这里obj的相关字段都要清洗掉，只留下必要的。
	obj.Status = NoUSE
	select {
	case t.pool <-obj:
		return nil
	default:
		return errors.New("put obj to pool error, the pool is full")
	}
}

func TestObjPool(t *testing.T) {

	pool := NewObjPool(5, false)

	obj1, err := pool.GetObj(time.Second * 10)
	obj1.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(obj1)

	for i := 0; i < 20; i++ {
		if obj, err := pool.GetObj(time.Second * 10); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println(obj)
			obj.Execute()
			fmt.Println(obj.GetId())
			if err := pool.PutObj(obj); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}
