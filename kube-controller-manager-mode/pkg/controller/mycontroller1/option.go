package mycontroller1

import "context"

type Option func(*ControllerOption)

// ControllerOption 控制器配置
type ControllerOption struct {
	ctx context.Context
	// queue 队列
	queue chan any
}

const (
	DefaultWorkers = 5
	miniWorkers    = 1
)

func WithContext(ctx context.Context) Option {
	return func(opt *ControllerOption) {
		opt.ctx = ctx
	}
}


func newOptions() *ControllerOption {
	return &ControllerOption{
		ctx:     context.Background(),
		queue:   make(chan interface{}, 10),
	}
}
