package kube_controller_mode

import (
	"fmt"
	"github.com/study-io/Golang-Concurrency-Pattern-Demo/gorountine-other-mode/group"
)

// Controller 控制器
type Controller struct {
	// ControllerOption 配置
	*ControllerOption
	// source 数据来源
	source <-chan interface{}
	// handler 当数据源进入时执行的方法
	handler ResourceHandler
	err     error
}

func NewController(opts ...Option) *Controller {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return &Controller{
		ControllerOption: options,
	}
}

func (c *Controller) Err() error {
	return c.err
}

// AddSource 设置控制器数据来源，
// 当控制器执行时，会不断从source chan中获取数据来源
func (c *Controller) AddSource(source <-chan interface{}) *Controller {
	c.source = source
	return c
}

// ResourceHandler 设置Handler
type ResourceHandler interface {
	SetHandler(obj interface{})
}

type ResourceHandlerFunc struct {
	SetHandlerFunc func(obj interface{})
}

// SetHandler 设置Handler
func (r ResourceHandlerFunc) SetHandler(obj interface{}) {
	if r.SetHandlerFunc != nil {
		r.SetHandlerFunc(obj)
	}
}

// AddHandler 设置控制器Handler，
// 当控制器执行时，会不断调用handler
func (c *Controller) AddHandler(handler ResourceHandler) *Controller {
	if c.source == nil {
		c.err = fmt.Errorf("controller source is nil")
		return c
	}

	c.handler = handler

	return c
}

// Run 执行控制器
func (c *Controller) Run() error {
	if c.err != nil {
		return c.err
	}

	if c.handler == nil {
		c.err = fmt.Errorf("controller handler is nil")
		return c.err
	}

	// 异步将来源的数据放入worker中
	go func() {
		for v := range c.source {
			c.queue <- v
		}
	}()

	// 并发启动worker
	g := group.NewGroup()
	for i := 0; i < c.workers; i++ {
		g.Start(c.runWorker)
	}
	g.Wait()

	return nil
}

// runWorker worker逻辑：不断从queue中获取
// 数据，并调用handler方法
func (c *Controller) runWorker() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		item, ok := <-c.queue
		if !ok && item == nil {
			return
		}
		c.handler.SetHandler(item)
	}
}
