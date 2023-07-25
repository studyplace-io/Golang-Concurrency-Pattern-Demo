package pipeline_mode

import (
	"fmt"
	"testing"
)


func TestPipeline(t *testing.T) {
	// 流水线模式
	in := producer(1, 2, 3, 4)
	ch := square(in)

	// 消费准备好的数据，ex: 落库等操作
	for ret := range ch {
		fmt.Println("res: ", ret)
	}

	fmt.Println("pipeline mode finished.")
}
