package forever_mode

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRunWithTimeForever(t *testing.T) {
	go func() {
		RunForeverWithTime(time.Second*1, func() error {
			fmt.Println("test-forever-with-time")
			return nil
		}, 3)
	}()

	<-time.After(time.Second * 10)

}

func TestRunWithTimeWithChannel(t *testing.T) {
	stopC := make(chan struct{})

	go func() {
		<-time.After(time.Second * 10)
		close(stopC)
	}()

	RunWithTimeWithChannel(time.Second*2, func() error {
		fmt.Println("test-with-time-channel")
		return nil
	}, stopC, 5)

}

func TestRunWithTimeWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		cancel()
	}()
	RunWithTimeWithContext(time.Second*1, func() error {
		fmt.Println("test-with-time-context")
		return nil
	}, ctx, 10)
}