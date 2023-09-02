package mycontroller2

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestController(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	controller, _ := NewController(WithContext(ctx), WithWorkers(5))

	ch := make(chan interface{}, 10)
	controller.AddSource(ch)

	go func() {
		for {
			ch <- rand.Intn(100)
			time.Sleep(time.Second * 1)
		}
	}()


	controller.AddHandler(ResourceHandlerFunc{
		SetHandlerFunc: func(obj interface{}) {
			fmt.Println("handler: ", obj)
		},
	})

	controller.Run()

}
