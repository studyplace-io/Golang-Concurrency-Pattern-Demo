package forever_mode

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRunForever(t *testing.T) {
	go func() {
		RunForever(time.Second*5, func() error {
			fmt.Println("test")
			return nil
		})
	}()

	<-time.After(time.Second * 20)

}

func TestRunWithChannel(t *testing.T) {
	stopC := make(chan struct{})

	go func() {
		<-time.After(time.Second * 20)
		close(stopC)
	}()

	RunWithChannel(time.Second*5, func() error {
		fmt.Println("test-channel")
		return nil
	}, stopC)

}

func TestRunWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer func() {
		cancel()
	}()
	RunWithContext(time.Second*5, func() error {
		fmt.Println("test-context")
		return nil
	}, ctx)
}
