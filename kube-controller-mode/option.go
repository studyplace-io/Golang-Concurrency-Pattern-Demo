package kube_controller_mode

import "context"

type Option func(*ControllerOption)

// ControllerOption 控制器配置
type ControllerOption struct {
	ctx context.Context
	// workers 并发worker数量
	workers int
	// queue 队列
	queue chan any
}

const (
	defaultWorkers = 5
	miniWorkers    = 1
)

func WithContext(ctx context.Context) Option {
	return func(opt *ControllerOption) {
		opt.ctx = ctx
	}
}

func WithWorkers(workers int) Option {
	return func(opts *ControllerOption) {
		if workers < 0 {
			opts.workers = miniWorkers
		} else {
			opts.workers = workers
		}
	}
}

func newOptions() *ControllerOption {
	return &ControllerOption{
		ctx:     context.Background(),
		workers: defaultWorkers,
		queue:   make(chan interface{}, 10),
	}
}
