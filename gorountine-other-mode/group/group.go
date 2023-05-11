package group

import (
	"context"
	"sync"
)

type Group struct {
	wg sync.WaitGroup
}

func NewGroup() *Group {
	return &Group{wg: sync.WaitGroup{}}
}

func (g *Group) Start(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		f()
	}()
}

func (g *Group) StartN(num int, f func()) {
	for i := 0; i < num; i++ {
		g.Start(f)
	}
}

func (g *Group) Wait() {
	g.wg.Wait()
}

func (g *Group) StartWithChannel(stopCh <-chan struct{}, f func(stopCh <-chan struct{})) {
	g.Start(func() {
		f(stopCh)
	})
}

func (g *Group) StartWithContext(ctx context.Context, f func(context.Context)) {
	g.Start(func() {
		f(ctx)
	})
}
