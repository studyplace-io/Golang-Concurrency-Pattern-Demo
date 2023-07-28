package worker_pool_mode

import (
	"fmt"
	"testing"
	"time"
)

func TestObjPool(t *testing.T) {

	pool := NewObjPool(5, false)

	obj1, err := pool.GetObj(time.Second * 5)
	obj1.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(obj1)

	for i := 0; i < 20; i++ {
		if obj, err := pool.GetObj(time.Second * 5); err != nil {
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
