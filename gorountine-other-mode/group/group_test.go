package group

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {

	g := NewGroup()

	g.Start(func() {
		time.Sleep(time.Second * 2)
		fmt.Println("test")
	})

	g.Wait()
}

func TestGroupN(t *testing.T) {

	g := NewGroup()

	g.StartN(5, func() {
		time.Sleep(time.Second * 2)
		fmt.Println("test")
	})

	g.Wait()
}

func TestGroupCtx(t *testing.T) {

	g := NewGroup()

	g.StartWithContext(context.TODO(), func(ctx context.Context) {
		time.Sleep(time.Second * 2)
		fmt.Println("test")
	})

	g.Wait()

}

func TestGroupStopC(t *testing.T) {

	g := NewGroup()
	stopC := make(chan struct{})
	g.StartWithChannel(stopC, func(stopCh <-chan struct{}) {
		time.Sleep(time.Second * 2)
		fmt.Println("test")
		select {
		case <-stopC:
			return
		}
	})
	time.Sleep(time.Second * 4)
	stopC <- struct{}{}

	g.Wait()

}
